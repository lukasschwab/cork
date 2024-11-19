package main

import (
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"syscall"

	"github.com/charmbracelet/log"
	"github.com/fsnotify/fsevents"
	"github.com/kballard/go-shellquote"
	"github.com/lukasschwab/cork/pkg/cork"
	"github.com/lukasschwab/cork/pkg/filter"
	"github.com/lukasschwab/cork/pkg/pattern"
)

// main kicks off the arg consumption cycle, then waits for an interrupt.
func main() {
	printConfigHeading()

	pairs := parse(os.Args[1:])
	watch(pairs)

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	<-sigChan
}

type configPair struct {
	command     string
	rawPatterns []string
	patterns    []pattern.Pattern
}

func (p configPair) String() string {
	return fmt.Sprintf("» %v → %s", p.rawPatterns, p.command)
}

func (p configPair) Action() cork.Action {
	action, err := execCommandAction(p.command)
	if err != nil {
		log.Fatalf("Error processing target command: %v", err)
	}

	filters := make([]filter.Func, len(p.patterns))
	for i := range p.patterns {
		filters[i] = p.patterns[i].Filter()
	}
	action.Filters = append(action.Filters, filter.Any(filters...))

	return action
}

func parse(args []string) (parsed []configPair) {
	cur := configPair{}
	for ; len(args) > 0; args = args[1:] {
		switch args[0] {
		// TODO: consider using filepath.SplitList and condensing multiple paths
		// to a single positional arg.
		case "-p", "--pattern":
			continue
		case "-r", "--run":
			// Ingest the command immediately.
			args = args[1:]
			cur.command = args[0]
			parsed = append(parsed, cur)
			cur = configPair{}
		default:
			if p, err := pattern.FromString(args[0]); err != nil {
				panic(err)
			} else {
				cur.rawPatterns = append(cur.rawPatterns, args[0])
				cur.patterns = append(cur.patterns, p)
			}
		}
	}

	return parsed
}

func watch(pairs []configPair) {
	patterns := []pattern.Pattern{}
	actions := make([]cork.Action, len(pairs))

	for i, pair := range pairs {
		printConfigPair(pair)
		patterns = append(patterns, pair.patterns...)
		actions[i] = pair.Action()
	}
	println()

	go func() {
		watcher := cork.Watcher{
			Paths:   patterns,
			Filters: []filter.Func{},
			Actions: actions,
		}
		if err := watcher.Watch(); err != nil {
			log.Fatalf("Failed to init watcher: %v", err)
		}
	}()
}

// TODO: consider a flag for enabling async exec; ruins output order.
func execCommandAction(cmdString string) (cork.Action, error) {
	splitCmd, err := shellquote.Split(cmdString)
	if err != nil {
		return cork.Action{}, err
	}

	callback := func(e fsevents.Event) {
		printEvent(e, cmdString)
		out, err := exec.Command(splitCmd[0], splitCmd[1:]...).Output()
		if err != nil {
			log.Error("Subprocess execution error", "message", err.Error())
		}
		if out != nil {
			// Bodge: make output visible.
			println(string(out))
		}
	}

	return cork.Action{Callback: callback}, nil
}
