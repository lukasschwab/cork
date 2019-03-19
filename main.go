package main

import (
	"log"
	"time"

	"./cork"
)

func main() {
	c, _ := cork.Init()
	defer c.Close()

	ag := cork.ActionGroup{
		Selector: func() []string {
			log.Println("From MAIN: selector running.")
			return []string{"./testdir"}
		},
		Filter: func(e cork.Event) bool {
			return true
		},
		Action: func(e cork.Event, cached string) string {
			return "Cache this, boys!"
		},
	}
	c.Add(ag)

	two := cork.ActionGroup{
		Selector: func() []string {
			log.Println("From MAIN: special selector running.")
			return []string{"./testdir/special"}
		},
		Filter: func(e cork.Event) bool {
			return true
		},
		Action: func(e cork.Event, cached string) string {
			return "Special cache."
		},
	}
	c.Add(two)

  // FIXME: run indefinitely.
	time.Sleep(300 * time.Second)
}
