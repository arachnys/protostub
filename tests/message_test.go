package tests

import (
	"fmt"
	"strings"
	"testing"

	"github.com/arachnys/protostub"
)

func expectMembers(t *testing.T, m *protostub.Message, expected map[string]bool) {
	for _, i := range m.Members {
		expected[fmt.Sprintf("%s %s", i.Typename(), i.Name())] = true
	}

	for k, v := range expected {
		if !v {
			t.Fatal("Not all members found. Failed to find", k)
		}
	}
}

func parse(t *testing.T, source string, types int) *protostub.ProtoData {
	p := protostub.New(strings.NewReader(source))

	err := p.Parse()

	if err != nil {
		t.Fatal(err)
	}

	if len(p.Types) != 1 {
		t.Fatal("Failed to parse proto ok")
	}

	return p
}

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

	p := parse(t, proto, 1)

	results := map[string]bool{
		"string thing":       false,
		"string other_thing": false,
		"Timestamp time":     false,
		"int number":         false,
		"Bar bar":            false,
	}

	if p.Types[0].Name() != "Foo" {
		t.Fatal("Failed to parse proto ok")
	}

	message, ok := p.Types[0].(*protostub.Message)

	if !ok {
		t.Fatal("Failed to read message ok")
	}

	expectMembers(t, message, results)
}

func TestMessageWithChildEnum(t *testing.T) {
	proto := `
	message Foo {
		enum Bar {
			A = 0;
			B = 1;
			C = 2;
		}
		Bar bar = 2;

		string name = 3;
	}
	`

	p := parse(t, proto, 1)

	message, ok := p.Types[0].(*protostub.Message)

	if !ok {
		t.Fatal("Failed to read message ok")
	}

	if len(message.Members) != 2 {
		t.Fatal("Failed to read message members")
	}

	if len(message.Types) != 1 {
		t.Fatal("Failed to read message subtypes")
	}

	fmt.Println(message.Types[0].Name())

	expected := map[string]bool{
		"A": false,
		"B": false,
		"C": false,
	}

	for _, i := range message.Types[0].(*protostub.Enum).Members {
		expected[i.Name()] = true
	}

	for _, v := range expected {
		if !v {
			t.Fatal("Failed to read enum properly")
		}
	}
}
