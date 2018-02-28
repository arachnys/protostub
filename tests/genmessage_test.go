package protostub

import (
	"bytes"
	"fmt"
	"strings"
	"testing"

	"github.com/arachnys/protostub"
	"github.com/arachnys/protostub/gen"
)

// This test is more of an experiment than anything
// using string comparison could make the test too brittle
// conversely, that could be exactly what is needed!

var sampleProto1 = `
message Foo {
    string bar = 1;
    int32 baz = 2;
}
`

// these are written without newlines or spaces
// while this makes the tests less comprehensive, it also makes them **far** less
// brittle and horrible to write.
var sampleMypy1 = "classFoo:bar:stringbaz:int"

var tests = [][]string{
	{sampleProto1, sampleMypy1},
}

func TestMessageGeneration(t *testing.T) {
	for _, i := range tests {
		proto := i[0]
		mypy := i[1]

		p := protostub.New(strings.NewReader(proto))

		err := p.Parse()

		if err != nil {
			t.Fatal(err)
		}

		buf := bytes.NewBuffer(nil)
		err = gen.Gen(buf, p)

		if err != nil {
			t.Fatal(err)
		}

		// remove newlines and spaces
		generated := strings.Replace(strings.Replace(buf.String(), " ", "", -1), "\n", "", -1)

		if generated != mypy {
			fmt.Println(generated)
			t.Fatal("Failed to generate correct code, got:\n", generated)
		}
	}
}
