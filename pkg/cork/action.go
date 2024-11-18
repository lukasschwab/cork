package cork

import (
	"github.com/fsnotify/fsnotify"
	"github.com/lukasschwab/cork/pkg/filter"
)

type Action struct {
	Filters  []filter.Func
	Callback func(e fsnotify.Event)
}

func (a *Action) handle(e fsnotify.Event) {
	if filter.All(a.Filters...)(e) {
		a.Callback(e)
	}
}
