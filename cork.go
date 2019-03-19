package cork

import (
	"crypto/md5"
	"fmt"
	"io"
	"log"
	"os"
	"sync"

	"github.com/fsnotify/fsnotify"
)

// An Event proxies an fsnotify.Event.
type Event struct {
	fsnotify.Event
}

// A Selector returns a list of relative file or directory names.
type Selector func() []string

// An Action receives an event and the previous cached value for the event file
// name. It returns the new value to be cached.
type Action func(e Event, cached string) string

// OnFileWrite runs A iff the event is a file write event.
func (a Action) OnFileWrite() Action {
	return func(e Event, cached string) string {
		if e.Op&fsnotify.Write == fsnotify.Write {
			return a(e, cached)
		}
		return cached
	}
}

// OnFileChange runs A iff the hash of the event file has changed. NOTE: this
// overrides A's cache values.
func (a Action) OnFileChange() Action {
	return func(e Event, cached string) string {
		newHash := fileHash(e.Name)
		if cached != newHash {
			a(e, cached)
		}
		return newHash
	}
}

// fileHash returns the md5 hash of file NAME.
func fileHash(name string) string {
	f, err := os.Open(name)
	if err != nil {
		log.Println("Failed to open file:", name)
		return ""
	}
	defer f.Close()

	h := md5.New()
	if _, err := io.Copy(h, f); err != nil {
		log.Println("Error generating hash for file:", name)
		return ""
	}
	return fmt.Sprintf("%x", h.Sum(nil))
}

// OnRegexChanges runs A iff the result of running regex R on the event file is
// novel.
func (a Action) OnRegexChanges() Action { // TODO
	return a
}

type Watcher struct {
	sync.RWMutex
	fsw   *fsnotify.Watcher
	cache map[string]string
}

func (w *Watcher) Close() {
	w.fsw.Close()
}

// getCache threadsafely retrieves the value associated with KEY in a watcher
// W's cache.
func (w *Watcher) getCache(key string) string {
	w.RLock()
	defer w.RUnlock()
	if val, ok := w.cache[key]; ok {
		return val
	}
	return ""
}

// setCache threadsafely sets the value associated with KEY in a watcher W's
// cache.
func (w *Watcher) setCache(key string, val string) {
	w.Lock()
	defer w.Unlock()
	w.cache[key] = val
}

// Watch creates returns a new Watcher that watches the files returned by S,
// and applies action A upon events from those files. You must call
// watcher.Close() to prevent memory leaks to fsnotify watchers.
//
// TODO: rerun selectors to find new files. Alternatively, depend on filter:
// watch all of the files in the cwd by default.
func Watch(s Selector, a Action) (*Watcher, error) {
	fsw, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, err
	}
	w := &Watcher{
		fsw:   fsw,
		cache: make(map[string]string),
	}

	go func() {
		for {
			select {
			case event, ok := <-w.fsw.Events:
				if !ok {
					log.Println("There was an error in an event consumer [events].")
					return
				}
				e := Event{event}
				cached := w.getCache(e.Name)
				w.setCache(e.Name, a(e, cached))
			case err, ok := <-w.fsw.Errors:
				if !ok {
					log.Println("There was an error in an event consumer [errs].")
					return
				}
				log.Println("error:", err)
			}
		}
	}()

	err = w.fsw.Add(s()[0])
	return w, err
}
