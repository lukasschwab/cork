package main

import (
	"./cork"
)

func main() {
	c, _ := cork.Init()
  defer c.Close()

  ag := cork.ActionGroup{
    Selector: func() []string {
      return []string{"testdir",}
    },
    Filter: func(e cork.Event) bool {
      return true
    },
    Action: func(e cork.Event) string {
      return "Cache this, boys!"
    },
  }
  c.Add(ag)
}
