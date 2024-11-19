package filter

import (
	"github.com/fsnotify/fsevents"
)

// TODO: consider exposing to the CLI.
func ByFlags(f fsevents.EventFlags) Func {
	return func(e fsevents.Event) bool {
		return e.Flags&f == f
	}
}
