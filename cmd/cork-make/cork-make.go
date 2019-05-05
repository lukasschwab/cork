package main

import (
	"log"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"strings"

	"github.com/fatih/color"
	"github.com/lukasschwab/cork"
)

// TODO: clean up logging.
var l = log.New(os.Stdout, "", 0)

// allWatchers is a global collection of watchers for cleanup.
var allWatchers = make([]*cork.Watcher, 0)

// Color-logging helpers.
var r = color.RedString
var g = color.GreenString
var b = color.BlueString

// main kicks off the arg consumption cycle, then waits for an interrupt.
func main() {
	defer cleanup()
	pwd, _ := filepath.Abs(".")
	l.Printf("Relative to %s:", pwd)
	parsePatterns(os.Args[1:])

	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt)
	select {
	case <-c:
		break
	}
}

// cleanup deprovisions watchers to prevent memory leaks.
func cleanup() {
	for _, w := range allWatchers {
		w.Close()
	}
}

// parsePatterns selects the leading patterns in the remaining args.
func parsePatterns(args []string) {
	if len(args) == 0 {
		println()
		return
	}
	if args[0] != "-p" && args[0] != "--pattern" {
		println(r("Error: first argument must be a pattern."))
		os.Exit(2)
	}
	var i int
	for i = 1; i < len(args) && args[i][0] != '-'; i++ {
	}
	parseCommand(args[1:i], args[i:])
}

// parseCommand selects the leading command from the remaining args.
func parseCommand(patterns []string, args []string) {
	if len(args) < 2 || (args[0] != "-r" && args[0] != "--run") {
		println(r("Error: patterns must be followed by a -r or --run."))
		os.Exit(2)
	}
	watch(patterns, args[1])
	parsePatterns(args[2:])
}

// watch spins up a watcher for the (PATTERNS, CMDSTRING) pair.
func watch(patterns []string, cmdString string) {
	println(g("» ['%s'] → %s", strings.Join(patterns, "', '"), cmdString))
	w, err := cork.Watch(selectPatterns(patterns), runCmd(cmdString).OnFileChange())
	if err != nil {
		println(r("Error creating watcher:"), err)
		return
	}
	allWatchers = append(allWatchers, w)
}

// selectPatterns returns the list of filenames that match the PATTERNS.
func selectPatterns(patterns []string) cork.Selector {
	return func() []string {
		var names = make(map[string]struct{})
		for _, p := range patterns {
			matches, _ := filepath.Glob(p) // FIXME: handle errors.
			for _, name := range matches {
				if _, in := names[name]; !in {
					names[name] = struct{}{}
				}
			}
		}
		unique := make([]string, len(names))
		i := 0
		for name := range names {
			unique[i] = name
			i++
		}
		return unique
	}
}

func runCmd(cmdString string) cork.Action {
	splitCmd := strings.Split(cmdString, " ") // FIXME: breaks on spaces in command args.
	return func(e cork.Event, cached string) string {
		println(b("%s → %s", e.Name, cmdString))
		out, err := exec.Command(splitCmd[0], splitCmd[1:]...).Output()
		if err != nil {
			println(r("Error:"), err)
		}
		if out != nil {
			println(b("output:"), string(out))
		}
		return "" // Discarded by OnFileChange().
	}
}
