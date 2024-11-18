package main

import (
	"log"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"

	"github.com/fatih/color"
	"github.com/fsnotify/fsnotify"
	"github.com/kballard/go-shellquote"
	"github.com/lukasschwab/cork/pkg/cork"
	"github.com/lukasschwab/cork/pkg/filter"
	"github.com/lukasschwab/cork/pkg/pattern"
)

var (
	NonChmod filter.Func = filter.Not(filter.ByOp(fsnotify.Chmod))
)

// TODO: clean up logging.
var l = log.New(os.Stdout, "", 0)

// TODO: use lipgloss
// Color-logging helpers.
var r = color.RedString
var g = color.GreenString
var b = color.BlueString

// main kicks off the arg consumption cycle, then waits for an interrupt.
func main() {
	pwd, _ := filepath.Abs(".")
	l.Printf("Relative to %s:", pwd)

	pairs := parse(os.Args[1:])
	watch(pairs)

	// Create a channel to receive interrupt signals
	sigChan := make(chan os.Signal, 1)

	// Notify the channel when an interrupt signal is received (e.g., Ctrl+C)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	// Block until a signal is received
	<-sigChan
}

type pair struct {
	command  string
	patterns []pattern.Pattern
}

func (p pair) Action() cork.Action {
	action, err := runCmd(p.command)
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

func parse(args []string) (parsed []pair) {
	cur := pair{}
	for ; len(args) > 0; args = args[1:] {
		switch args[0] {
		case "-p", "--pattern":
			continue
		case "-r", "--run":
			// Ingest the command at the same time.
			args = args[1:]
			cur.command = args[0]
			parsed = append(parsed, cur)
			cur = pair{}
		default:
			if p, err := pattern.FromString(args[0]); err != nil {
				panic(err)
			} else {
				cur.patterns = append(cur.patterns, p)
			}
		}
	}

	return parsed
}

func watch(pairs []pair) {
	patterns := []pattern.Pattern{}
	actions := make([]cork.Action, len(pairs))

	for i, pair := range pairs {
		patternStrings := make([]string, 0, len(pair.patterns))
		for _, pat := range pair.patterns {
			patternStrings = append(patternStrings, string(pat))
		}
		println(g("» ['%s'] → %s", strings.Join(patternStrings, "', '"), pair.command))

		patterns = append(patterns, pair.patterns...)
		actions[i] = pair.Action()
	}

	go func() {
		watcher := cork.Watcher{
			Paths:   patterns,
			Filters: []filter.Func{NonChmod},
			Actions: actions,
		}
		watcher.Watch()
	}()
}

func runCmd(cmdString string) (cork.Action, error) {
	splitCmd, err := shellquote.Split(cmdString)
	if err != nil {
		return cork.Action{}, err
	}

	cb := func(e fsnotify.Event) {
		println(b("%s → %s", e.Name, cmdString))
		out, err := exec.Command(splitCmd[0], splitCmd[1:]...).Output()
		if err != nil {
			println(r("Error:"), err.Error())
		}
		if out != nil {
			println(b("output:"), string(out))
		}
	}

	return cork.Action{Callback: cb}, nil
}
