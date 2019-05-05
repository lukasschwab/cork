package main

import (
  "log"
  "os"
  "strings"

  // "github.com/lukasschwab/cork"
)

var stderr = log.New(os.Stderr, "", 0)
var stdout = log.New(os.Stderr, "", 0)

func main() {
  parsePatterns(os.Args[1:])
}

func parsePatterns(args []string) {
  if len(args) == 0 {
    return
  }
  if args[0] != "-p" && args[0] != "--pattern" {
    stderr.Println("First argument must be a pattern.")
    os.Exit(2)
  }
  var i int;
  for i = 1; i < len(args) && args[i][0] != '-'; i++ {}
  parseCommand(args[1:i], args[i:])
}

func parseCommand(patterns []string, args []string) {
  if len(args) < 2 || (args[0] != "-r" && args[0] != "--run") {
    stderr.Println("Patterns must be followed by a -r or --run.")
    os.Exit(2)
  }
  command := args[1]
  stdout.Printf("['%s'] → %s", strings.Join(patterns, "', '"), command)
  parsePatterns(args[2:])
}

func watch(patterns []string, cmd string) {
  // TODO: set up watcher.
}
