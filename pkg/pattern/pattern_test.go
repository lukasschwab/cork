package pattern_test

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/lukasschwab/cork/pkg/pattern"
	"github.com/peterldowns/testy/assert"
)

func TestGetStaticPrefix(t *testing.T) {
	wd, err := os.Getwd()
	assert.Nil(t, err)

	cases := map[string]string{
		// ABSOLUTE CASES
		// Absolute, no wildcard
		"/":      "/",
		"/inter": "/inter",
		// Absolute wildcard
		"/*":       "/",
		"/pre*":    "/",
		"/pre*suf": "/",
		"/*suf":    "/",
		// Absolute wildcard with intermediates
		"/inter/*":       "/inter",
		"/inter/pre*":    "/inter",
		"/inter/pre*suf": "/inter",
		"/inter/*suf":    "/inter",

		// Local, no wildcard
		"inter": filepath.Join(wd, "inter"),
		// Local wildcard
		"./*":       wd,
		"*":         wd,
		"./pre*":    wd,
		"pre*":      wd,
		"./pre*suf": wd,
		"pre*suf":   wd,
		"./*suf":    wd,
		"*suf":      wd,
		// Local wildcard with intermediates
		"./inter/*":       filepath.Join(wd, "inter"),
		"inter/*":         filepath.Join(wd, "inter"),
		"./inter/pre*":    filepath.Join(wd, "inter"),
		"inter/pre*":      filepath.Join(wd, "inter"),
		"./inter/pre*suf": filepath.Join(wd, "inter"),
		"inter/pre*suf":   filepath.Join(wd, "inter"),
		"./inter/*suf":    filepath.Join(wd, "inter"),
		"inter/*suf":      filepath.Join(wd, "inter"),
	}

	for input, expected := range cases {
		t.Run(input, func(t *testing.T) {
			p, err := pattern.FromString(input)
			assert.Nil(t, err)
			res := p.StaticPrefixPath()
			assert.Equal(t, expected, res)
		})
	}
}

func TestCanMatchChildrenOf(t *testing.T) {
	cases := map[string]string{
		"a/b/c/**/*.txt":   "a/b/c/1",
		"a/b/c/*/*.txt":    "a/b/c/1",
		"./a/b/c/**/*.txt": "a/b/c/1",
		"./a/b/c/*/*.txt":  "a/b/c/1",
		"/a/b/c/**/*.txt":  "/a/b/c/1",
		"/a/b/c/*/*.txt":   "/a/b/c/1",
	}

	for input, path := range cases {
		t.Run(fmt.Sprintf("%s should match %s children", input, path), func(t *testing.T) {
			p, _ := pattern.FromString(input)
			canMatch, err := p.CanMatchChildrenOf(path)
			assert.Nil(t, err)
			assert.True(t, canMatch)
		})
	}
}
