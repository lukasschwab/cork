package main

import (
	"log"
	"os"
	"os/exec"
	"os/signal"
	"strings"

	"github.com/fatih/color"

	"github.com/lukasschwab/cork"
)

var stderr = log.New(os.Stderr, "", 0)
var stdout = log.New(os.Stderr, "", 0)

var allWatchers = make([]*cork.Watcher, 0)

var r = color.RedString
var g = color.GreenString
var b = color.BlueString

func main() {
	defer cleanup()
	parsePatterns(os.Args[1:])

	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt)
	select {
	case <-c:
		break
	}
}

func cleanup() {
	for _, w := range allWatchers {
		w.Close()
	}
}

func parsePatterns(args []string) {
	if len(args) == 0 {
		println()
		return
	}
	if args[0] != "-p" && args[0] != "--pattern" {
		stderr.Println("First argument must be a pattern.")
		os.Exit(2)
	}
	var i int
	for i = 1; i < len(args) && args[i][0] != '-'; i++ {
	}
	parseCommand(args[1:i], args[i:])
}

func parseCommand(patterns []string, args []string) {
	if len(args) < 2 || (args[0] != "-r" && args[0] != "--run") {
		stderr.Println("Patterns must be followed by a -r or --run.")
		os.Exit(2)
	}
	watch(patterns, args[1])
	parsePatterns(args[2:])
}

func watch(patterns []string, cmdString string) {
	stdout.Print(g("» ['%s'] → %s", strings.Join(patterns, "', '"), cmdString))
	w, err := cork.Watch(cork.SelectPatterns(patterns), runCmd(cmdString).OnFileChange())
	if err != nil {
		stderr.Println("Error creating watcher:", err)
		return
	}
	allWatchers = append(allWatchers, w)
}

func runCmd(cmdString string) cork.Action {
	splitCmd := strings.Split(cmdString, " ") // FIXME: breaks on spaces in command args.
	return func(e cork.Event, cached string) string {
		stdout.Print(b("%s → %s", e.Name, cmdString))
		out, err := exec.Command(splitCmd[0], splitCmd[1:]...).Output()
		if err != nil {
			stderr.Println(r("Error:"), err)
		}
		if out != nil {
			stdout.Println(b("output:"), string(out))
		}
		return "" // Discarded by OnFileChange().
	}
}
