package main

import (
	"os"
	"os/exec"
	"time"

	"./cork"
)

func main() {
	c, _ := cork.Init()
	defer c.Close()

	ag := cork.ActionGroup{
		Selector: func() []string {
			return []string{"./testdir"}
		},
		Filter: func(e cork.Event, cached string) bool {
			return true
		},
		Action: func(e cork.Event, cached string) string {
			return cached
		},
	}
	c.Add(ag)

	two, _ := cork.ActWhenFileChanges(cork.ActionGroup{
		Selector: func() []string {
			return []string{"./testdir"}
		},
		Filter: func(e cork.Event, cached string) bool {
			return true
		},
		Action: func(e cork.Event, cached string) string {
			cmd := exec.Command("ls", ".").Output()
			cmd.Stdout = os.Stdout
			return "" // Return value is discarded.
		},
	})
	c.Add(two)

	// FIXME: run indefinitely.
	time.Sleep(300 * time.Second)
}
