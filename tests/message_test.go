package protostub

import (
	"fmt"
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

func TestMessageMembers(t *testing.T) {
	proto := `
	message Foo {
		string thing = 1;
		string other_thing = 2;
		google.protobuf.Timestamp time = 3;
		uint64 number = 4;
		Bar bar = 5;
	}
	`

	results := map[string]bool{
		"string thing":       false,
		"string other_thing": false,
		"Timestamp time":     false,
		"uint64 number":      false,
		"Bar bar":            false,
	}

	p := protostub.New(strings.NewReader(proto))

	err := p.Parse()

	if err != nil {
		t.Fatal(err)
	}

	if len(p.Types) != 1 {
		t.Fatal("Failed to parse proto ok")
	}

	if p.Types[0].Name() != "Foo" {
		t.Fatal("Failed to parse proto ok")
	}

	message, ok := p.Types[0].(*protostub.Message)

	if !ok {
		t.Fatal("Failed to read message ok")
	}

	for _, i := range message.Members {
		results[fmt.Sprintf("%s %s", i.Typename(), i.Name())] = true
	}

	for _, v := range results {
		if !v {
			t.Fatal("Not all members found", results)
		}
	}
}
