package protostub

import (
	"strings"
	"testing"
)

// Just test that a message can be visited and the name is read all OK
func TestVisitMessage(t *testing.T) {
	proto := [][]string{
		{"message Foo {}", "Foo"},
		{"message Bar {}", "Bar"},
		{"message Baz {}", "Baz"},
	}

	for _, i := range proto {
		v := New(strings.NewReader(i[0]))
		err := v.Parse()

		if err != nil {
			t.Fatal(err)
		}
	}
}
