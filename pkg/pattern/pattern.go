package pattern

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/fsnotify/fsevents"
	"github.com/lukasschwab/cork/pkg/filter"
)

type Pattern string

func FromString(p string) (Pattern, error) {
	abs, err := filepath.Abs(string(p))
	if err != nil {
		return "", fmt.Errorf("can't make pattern absolute: %w", err)
	} else if _, err := filepath.Match(abs, ""); err != nil {
		return "", fmt.Errorf("absolute pattern isn't usable: %w", err)
	}
	return Pattern(abs), nil
}

func (p Pattern) Filter() filter.Func {
	pattern, err := filepath.Abs(string(p))
	if err != nil {
		panic(err)
	} else if _, err := filepath.Match(pattern, ""); err != nil {
		panic(fmt.Errorf("invalid pattern: %w", err))
	}

	return func(e fsevents.Event) bool {
		// Pattern is valid: checked above.
		match, _ := filepath.Match(pattern, e.Path)
		return match
	}
}

// StaticPrefixPath makes a best effort at getting a low, non-wildcard path to
// watch. The lower it is, the fewer superfluous events it'll receive.
func (p Pattern) StaticPrefixPath() string {
	// Drop anything after a `*` wildcard.
	if i := strings.Index(string(p), "*"); i != -1 {
		prefix, _ := filepath.Split(string(p)[:i])
		return filepath.Clean(prefix)
	}
	return string(p)
}
