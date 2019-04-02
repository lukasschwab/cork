package cork

import (
	"log"
	"testing"
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

	// Run indefinitely.
	var c chan struct{}
	<-c
}

func TestCork(t *testing.T) {
	return
}
