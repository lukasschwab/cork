package main

import (
  "log"
	"time"
  // "os/exec"

	"./cork"
)

func main() {
	c, _ := cork.Init()
	defer c.Close()

  var normalAction cork.Action = func(e cork.Event, cached string) string {
    log.Println("Normal cache:", cached)
    return "The normal cache never changes."
  }

	// ag := cork.ActionGroup{
	// 	Selector: func() []string {
	// 		return []string{"./testdir"}
	// 	},
	// 	Action: normalAction.OnFileWrite(),
	// }
	c.Add(func() []string {
    return []string{"./testdir"}
  }, normalAction.OnFileWrite())

  var specialAction cork.Action = func(e cork.Event, cached string) string {
    log.Println("Filechange cache:", cached)
    return ""
  }

	// two := cork.ActionGroup{
  //   Selector: func() []string {
	// 		return []string{"./testdir"}
	// 	},
	// 	Action: specialAction.OnFileChange().OnFileWrite(),
  // }
	c.Add(func() []string {
    return []string{"./testdir"}
  }, specialAction.OnFileChange().OnFileWrite())

  // FIXME: run indefinitely.
	time.Sleep(300 * time.Second)
}
