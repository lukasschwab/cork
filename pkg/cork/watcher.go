// Package cork is a file event handler.
package cork

import (
	"fmt"
	"time"

	"github.com/fsnotify/fsevents"
	"github.com/lukasschwab/cork/pkg/filter"
	"github.com/lukasschwab/cork/pkg/pattern"
)

const (
	// DefaultLatency for filesystem events before processing by Watcher.
	// fsevents docs suggest this helps debounce.
	DefaultLatency = 200 * time.Millisecond
)

// Watcher for filesystem events.
type Watcher struct {
	Paths   []pattern.Pattern
	Filters []filter.Func
	Actions []Action

	stream *fsevents.EventStream
}

// Watch events on w.Paths, filter with w.Filters, and trigger w.Actions in a
// spawned goroutine.
func (w *Watcher) Watch() error {
	pathsToWatch := make([]string, len(w.Paths))
	for i, path := range w.Paths {
		pathsToWatch[i] = path.StaticPrefixPath()
	}

	w.stream = &fsevents.EventStream{
		Paths:   pathsToWatch,
		Latency: DefaultLatency,
		Flags:   fsevents.FileEvents | fsevents.IgnoreSelf,
	}

	if err := w.stream.Start(); err != nil {
		return fmt.Errorf("can't start stream: %w", err)
	}

	go w.loop()
	return nil
}

func (w *Watcher) loop() {
	for eventGroup := range w.stream.Events {
		// Handle a file at most once per batch of events.
		handled := map[string]bool{}
		for _, event := range eventGroup {
			if _, ok := handled[event.Path]; !ok {
				w.handle(event)
				handled[event.Path] = true
			}
		}
	}
}

func (w *Watcher) handle(e fsevents.Event) {
	if filter.All(w.Filters...)(e) {
		for _, a := range w.Actions {
			a.handle(e)
		}
	}
}
