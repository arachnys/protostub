package protostub

import (
	"strings"
	"testing"

	"github.com/arachnys/protostub"
)

// Just test that a message can be visited and the name is read all OK
func TestVisitMessage(t *testing.T) {
	proto := [][]string{
		{"message Foo {}", "Foo"},
		{"message Bar {}", "Bar"},
		{"message Baz {}", "Baz"},
	}

	for _, i := range proto {
		v := protostub.New(strings.NewReader(i[0]))
		err := v.Parse()

		if err != nil {
			t.Fatal(err)
		}

		if len(v.Types) != 1 {
			t.Fatal("Failed to parse message")
		}

		if v.Types[0].Name() != i[1] {
			t.Fatal("Failed to read message name")
		}
	}
}
