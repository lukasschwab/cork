package pattern

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/fsnotify/fsnotify"
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

	return func(e fsnotify.Event) bool {
		// Pattern is valid: checked above.
		match, _ := filepath.Match(pattern, e.Name)
		return match
	}
}

func (p Pattern) StaticPrefixPath() string {
	// Drop anything after a `*` wildcard.
	if i := strings.Index(string(p), "*"); i != -1 {
		prefix, _ := filepath.Split(string(p)[:i])
		return filepath.Clean(prefix)
	}
	return string(p)
}

func (p Pattern) WildcardContainingDirectories() []string {
	prefixesToWildcards := []string{string(p)}
	for i, c := range p {
		if c == '*' {
			prefix, _ := filepath.Split(string(p)[:i])
			prefixesToWildcards = append(prefixesToWildcards, filepath.Clean(prefix))
		}
	}
	return prefixesToWildcards
}

func (p Pattern) CanMatchChildrenOf(path string) (match bool, err error) {
	if path, err = filepath.Abs(path); err != nil {
		return false, fmt.Errorf("error making target path absolute: %w", err)
	}

	for prefix := string(p); prefix != "/"; prefix = filepath.Dir(prefix) {
		if match, _ = filepath.Match(prefix, path); match {
			return true, nil
		}
	}
	return false, nil
}
