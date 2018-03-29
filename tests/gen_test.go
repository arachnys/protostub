package tests

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

var simpleProto = `
message Foo {
	string bar = 1;
	int32 baz = 2;
}
`

// these are written without newlines or spaces
// while this makes the tests less comprehensive, it also makes them **far** less
// brittle and horrible to write.
var simpleMypy = "fromtypingimportcast,Dict,List,TupleclassFoo:bar:stringbaz:intdef__init__(self,bar:string=None,baz:int=None)->Foo:...defCopyFrom(self,other:Foo)->None:...@staticmethoddefListFields()->Tuple[FieldDescriptor,value]:..."

var enumProto = `
enum Foo {
	BAR = 0;
	BAZ = 1;
	QUUX = 5;
}
`

var enumMypy = "fromtypingimportcast,Dict,List,TupleclassFoo:BAR=cast(Foo,0)BAZ=cast(Foo,1)QUUX=cast(Foo,5)@staticmethoddefName(number:Foo)->str:...@staticmethoddefValue(name:str)->Foo:...defCopyFrom(self,other:Foo)->None:...@staticmethoddefListFields()->Tuple[FieldDescriptor,value]:..."

var tests = [][]string{
	{simpleProto, simpleMypy},
	{enumProto, enumMypy},
}

func TestGeneration(t *testing.T) {
	for _, i := range tests {
		proto := i[0]
		mypy := i[1]

		p := protostub.New(strings.NewReader(proto))

		if err := p.Parse(); err != nil {
			t.Fatal(err)
		}

		buf := bytes.NewBuffer(nil)

		if err := gen.Gen(buf, p, false); err != nil {
			t.Fatal(err)
		}

		// remove newlines and spaces
		generated := strings.Replace(strings.Replace(buf.String(), " ", "", -1), "\n", "", -1)

		if generated != mypy {
			fmt.Println(mypy)
			fmt.Println(generated)
			t.Fatal("Failed to generate correct code, got:\n", generated)
		}
	}
}
