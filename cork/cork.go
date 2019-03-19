package cork

import (
	"log"
	"sync"

	"github.com/fsnotify/fsnotify"
)

type Event struct {
	fsnotify.Event
}

type FileSelector func() []string

type TriggerFilter func(e Event) bool

type ActionIntoCache func(e Event, cached string) string

type ActionGroup struct {
	Selector FileSelector
	Filter   TriggerFilter
	Action   ActionIntoCache
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
				if ag.Filter(e) {
					newVal := ag.Action(e, c.GetCache(e.Name))
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
