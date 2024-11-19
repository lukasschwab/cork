package pattern_test

import (
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
