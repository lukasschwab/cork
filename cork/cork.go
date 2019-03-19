package cork

import (
	"log"

	"github.com/fsnotify/fsnotify"
	"github.com/patrickmn/go-cache"
)

type Event struct {
  fsnotify.Event
}

type FileSelector func() []string

type TriggerFilter func(e Event) bool

type ActionIntoCache func(e Event) string

type ActionGroup struct {
	Selector FileSelector
	Filter   TriggerFilter
	Action   ActionIntoCache
}

type Cork struct {
	cache *cache.Cache
  watchers []*fsnotify.Watcher
}

func Init() (*Cork, error) {
	return &Cork{
		cache: cache.New(cache.NoExpiration, cache.NoExpiration),
    watchers: make([]*fsnotify.Watcher, 0),
	}, nil
}

func (c *Cork) Close() {
  for _, w := range c.watchers {
    log.Println("Destroying a watcher.")
    w.Close()
  }
}

func (c *Cork) Add(ag ActionGroup) error {
  w, err := fsnotify.NewWatcher()
  if err != nil {
    return err
  }
  c.watchers = append(c.watchers, w)

  return nil
}
