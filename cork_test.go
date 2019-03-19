package cork

import (
	"log"
	"testing"
	"time"
)

func Example() {
	var normalAction Action = func(e Event, cached string) string {
		log.Println("Normal cache:", cached)
		return "The normal cache never changes."
	}

	w1, _ := Watch(func() []string {
		return []string{"./testdir"}
	}, normalAction.OnFileWrite())
	defer w1.Close()

	var specialAction Action = func(e Event, cached string) string {
		log.Println("Filechange cache:", cached)
		return ""
	}

	w2, _ := Watch(func() []string {
		return []string{"./testdir"}
	}, specialAction.OnFileChange().OnFileWrite())
	defer w2.Close()

	// FIXME: run indefinitely.
	time.Sleep(300 * time.Second)
}

func TestCork(t *testing.T) {
	return
}
