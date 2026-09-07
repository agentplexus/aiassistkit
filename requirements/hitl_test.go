package requirements

import (
	"bytes"
	"strings"
	"testing"
)

// TestCLIPrompterMultipleCallsOverBufferedInput is a regression test: a
// bufio.Reader created fresh per call can read ahead into a pipe carrying
// several answers at once, silently dropping whatever was buffered but unread
// when a new reader is created for the next call. The reader must persist
// across calls.
func TestCLIPrompterMultipleCallsOverBufferedInput(t *testing.T) {
	in := strings.NewReader("2\ny\n1\n")
	var out bytes.Buffer
	p := &CLIPrompter{In: in, Out: &out}

	choice, err := p.Choose("pick one", []string{"a", "b", "c"})
	if err != nil {
		t.Fatalf("first Choose: %v", err)
	}
	if choice != 1 {
		t.Fatalf("first Choose = %d, want 1", choice)
	}

	confirmed, err := p.Confirm("proceed?")
	if err != nil {
		t.Fatalf("Confirm: %v", err)
	}
	if !confirmed {
		t.Fatalf("Confirm = false, want true")
	}

	choice2, err := p.Choose("pick one", []string{"x", "y", "z"})
	if err != nil {
		t.Fatalf("second Choose: %v", err)
	}
	if choice2 != 0 {
		t.Fatalf("second Choose = %d, want 0", choice2)
	}
}

// TestCLIPrompterStructLiteralWorks confirms the lazy reader init means a
// CLIPrompter built without NewCLIPrompter (e.g. &CLIPrompter{In: r, Out: w})
// still behaves correctly.
func TestCLIPrompterStructLiteralWorks(t *testing.T) {
	p := &CLIPrompter{In: strings.NewReader("yes\n"), Out: &bytes.Buffer{}}
	ok, err := p.Confirm("?")
	if err != nil {
		t.Fatalf("Confirm: %v", err)
	}
	if !ok {
		t.Fatalf("Confirm = false, want true")
	}
}
