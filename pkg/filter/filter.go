package filter

import (
	"github.com/fsnotify/fsevents"
)

var (
	Trivial Func = func(_ fsevents.Event) bool {
		return true
	}
)

// TODO: consider conversion to an interface.
type Func func(fsevents.Event) bool

func (f Func) Or(other Func) Func {
	return func(e fsevents.Event) bool {
		return f(e) || other(e)
	}
}

func (f Func) And(other Func) Func {
	return func(e fsevents.Event) bool {
		return f(e) && other(e)
	}
}

func Not(filter Func) Func {
	return func(e fsevents.Event) bool {
		return !filter(e)
	}
}

func All(filters ...Func) Func {
	return func(e fsevents.Event) bool {
		for _, f := range filters {
			if !f(e) {
				return false
			}
		}
		return true
	}
}

func Any(filters ...Func) Func {
	return func(e fsevents.Event) bool {
		for _, f := range filters {
			if f(e) {
				return true
			}
		}
		return false
	}
}
