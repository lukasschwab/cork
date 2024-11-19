package cork

import (
	"github.com/fsnotify/fsevents"
	"github.com/lukasschwab/cork/pkg/filter"
)

type Action struct {
	Filters  []filter.Func
	Callback func(e fsevents.Event)
}

func (a *Action) handle(e fsevents.Event) {
	if filter.All(a.Filters...)(e) {
		a.Callback(e)
	}
}
