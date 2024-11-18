package filter

import (
	"github.com/fsnotify/fsnotify"
)

func ByOp(op fsnotify.Op) Func {
	return func(e fsnotify.Event) bool {
		return e.Op.Has(op)
	}
}
