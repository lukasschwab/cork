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

type Event struct {
	fsnotify.Event
}

type Selector func() []string

type Action func(e Event, cached string) string

func (a Action) OnFileChange() Action {
  return func(e Event, cached string) string {
    newHash := fileHash(e.Name)
    if cached != newHash {
      a(e, cached)
    }
    return newHash
  }
}

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

func (a Action) OnFileWrite() Action {
  return func(e Event, cached string) string {
    if e.Op&fsnotify.Write == fsnotify.Write {
      return a(e, cached)
    }
    return cached
  }
}

func (a Action) OnRegexChanges() Action {
  return a
}

type watcher struct {
  sync.RWMutex
  fsw *fsnotify.Watcher
  cache map[string]string
}

func (w *watcher) close() {
  w.fsw.Close()
}

func (w *watcher) GetCache(key string) string {
	w.RLock()
	defer w.RUnlock()
	if val, ok := w.cache[key]; ok {
		return val
	}
	return ""
}

func (w *watcher) SetCache(key string, val string) {
	w.Lock()
	defer w.Unlock()
	w.cache[key] = val
}

type Cork struct {
	sync.RWMutex
	cache    map[string]string
	watchers []*watcher
}

func Init() (*Cork, error) {
	return &Cork{
		watchers: make([]*watcher, 0),
	}, nil
}

func (c *Cork) Close() {
	for _, w := range c.watchers {
		log.Println("Destroying a watcher.")
		w.close()
	}
}

// TODO: rerun selectors to find new files. Alternatively, depend on filter:
// watch all of the files in the cwd by default.
func (c *Cork) Add(s Selector, a Action) error {
	fsw, err := fsnotify.NewWatcher()
	if err != nil {
		return err
	}
  w := &watcher{
    fsw: fsw,
    cache: make(map[string]string),
  }

	c.watchers = append(c.watchers, w)

	go func() {
		for {
			select {
			case event, ok := <-w.fsw.Events:
				if !ok {
					log.Println("There was an error in an event consumer [events].")
					return
				}
				e := Event{event}
				cached := w.GetCache(e.Name)
				w.SetCache(e.Name, a(e, cached))
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
	return err
}
