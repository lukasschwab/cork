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
      return []string{"./testdir",}
    },
    Filter: func(e cork.Event) bool {
      log.Println("From MAIN: filter running.")
      return true
    },
    Action: func(e cork.Event) string {
      log.Println("From MAIN: action running.")
      return "Cache this, boys!"
    },
  }
  c.Add(ag)

  time.Sleep(300 * time.Second)
}
