package cork

import (
	"crypto/md5"
	"errors"
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

type FileSelector func() []string

type TriggerFilter func(e Event, cached string) bool

type ActionToCache func(e Event, cached string) string

type ActionGroup struct {
	Selector FileSelector
	Filter   TriggerFilter
	Action   ActionToCache
}

type Cork struct {
	sync.RWMutex
	cache    map[string]string
	watchers []*fsnotify.Watcher
}

func Init() (*Cork, error) {
	return &Cork{
		cache:    make(map[string]string),
		watchers: make([]*fsnotify.Watcher, 0),
	}, nil
}

func (c *Cork) GetCache(key string) string {
	c.RLock()
	defer c.RUnlock()
	if val, ok := c.cache[key]; ok {
		return val
	}
	return ""
}

func (c *Cork) SetCache(key string, val string) {
	c.Lock()
	defer c.Unlock()
	c.cache[key] = val
}

func (c *Cork) Close() {
	for _, w := range c.watchers {
		log.Println("Destroying a watcher.")
		w.Close()
	}
}

// TODO: rerun selectors to find new files. Alternatively, depend on filter:
// watch all of the files in the cwd by default.
// TODO: should the cache exist on a group level or on a global level?
// Probably on a group level...
func (c *Cork) Add(ag ActionGroup) error {
	w, err := fsnotify.NewWatcher()
	if err != nil {
		return err
	}
	c.watchers = append(c.watchers, w)

	go func() {
		for {
			select {
			case event, ok := <-w.Events:
				if !ok {
					log.Println("There was an error in an event consumer [events].")
					return
				}
				e := Event{event}
				cached := c.GetCache(e.Name)
				if ag.Filter(e, cached) {
					newVal := ag.Action(e, cached)
					c.SetCache(e.Name, newVal)
				}
			case err, ok := <-w.Errors:
				if !ok {
					log.Println("There was an error in an event consumer [errs].")
					return
				}
				log.Println("error:", err)
			}
		}
	}()

	err = w.Add(ag.Selector()[0])

	return nil
}

func ActWhenFileChanges(ag ActionGroup) (ActionGroup, error) {
	if ag.Selector == nil {
		return ActionGroup{}, errors.New("No FileSelector defined.")
	}
	return ActionGroup{
		Selector: ag.Selector,
		Filter: func(e Event, cached string) bool {
			return ag.Filter(e, cached) && cached != fileHash(e.Name)
		},
		Action: func(e Event, cached string) string {
			if ag.Action != nil {
				ag.Action(e, cached)
			}
			return fileHash(e.Name) // TODO: get the real hash.
		},
	}, nil
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

// TODO: uses regex to filter for a subset of the file to check changes on.
// Useful: https://golang.org/pkg/regexp/#Regexp.FindAllString
func ActWhenRegexChanges() {
	return
}
