// Package cork is a file event handler.
package cork

import (
	"log"
	"os"
	"path/filepath"

	"github.com/fsnotify/fsnotify"
	"github.com/lukasschwab/cork/pkg/filter"
	"github.com/lukasschwab/cork/pkg/pattern"
)

// TODO: consider adding support for on-the-fly addition of paths, filters,
// actions. There's support in fsnotify, but I think it adds too much complexity
// (async registration) for an application where the config is known.
//
// Basically, this should be driven by the need for an async Watch function. I
// don't actually have a need for that right now.
type Watcher struct {
	Paths   []pattern.Pattern
	Filters []filter.Func
	Actions []Action
}

func (w *Watcher) handle(e fsnotify.Event) {
	if filter.All(w.Filters...)(e) {
		// Note: this is sequential rather than parallel for actions.
		for _, a := range w.Actions {
			a.handle(e)
		}
	}
}

func (w *Watcher) Watch() {
	inner, err := fsnotify.NewWatcher()
	if err != nil {
		panic(err)
	}

	// Block indefinitely and handle events.
	go func() {
		for event := range inner.Events {
			w.handle(event)

			// Check if the event is for a directory
			if info, err := os.Lstat(event.Name); event.Has(fsnotify.Create) && err == nil && info.IsDir() {
				// Check if can prefix
				isRelevant := false
				for _, p := range w.Paths {
					if match, _ := p.CanMatchChildrenOf(event.Name); match {
						isRelevant = true
						break
					}
				}
				if isRelevant {
					w.addSubdirectory(event.Name, inner)
				}
			}
		}
	}()

	// TODO: discover glob-match directories upon initialization.
	for _, p := range w.Paths {
		for _, pathPrefix := range p.WildcardContainingDirectories() {
			matches, err := filepath.Glob(pathPrefix)
			if err != nil {
				log.Panicf("Couldn't glob prefix %s", pathPrefix)
			}
			for _, match := range matches {
				if err := inner.Add(match); err != nil {
					log.Panicf("Error watching path: %v", err)
				}
			}
		}
	}
}

func (w *Watcher) addSubdirectory(path string, inner *fsnotify.Watcher) {
	if newPattern, err := pattern.FromString(path); err != nil {
		log.Printf("Can't watch new subdirectory: %v", err)
	} else {
		log.Printf("Watching new subdirectory %s", path)
		w.Paths = append(w.Paths, newPattern)
		if err := inner.Add(string(newPattern)); err != nil {
			log.Printf("Can't watch new subdirectory: %v", err)
		}
	}
}
