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
var simpleMypy = "fromtypingimportcast,Dict,List,TupleclassFoo:bar:strbaz:intdef__init__(self,bar:str=None,baz:int=None)->Foo:...defCopyFrom(self,other:Foo)->None:...defSerializeToString(self)->str:...defParseFromString(self,data:str)->None:...@staticmethoddefListFields()->Tuple[FieldDescriptor,value]:..."

var enumProto = `
enum Foo {
	BAR = 0;
	BAZ = 1;
	QUUX = 2;
}
`

var enumMypy = "fromtypingimportcast,Dict,List,TupleclassFoo:BAR:AnyBAZ:AnyQUUX:Anydef__init__(self,BAR:Any=None,BAZ:Any=None,QUUX:Any=None)->Foo:...defName(enumClass:Foo)->Any:...defValue(memberName:str)->Any:...defCopyFrom(self,other:Foo)->None:...defSerializeToString(self)->str:...defParseFromString(self,data:str)->None:...@staticmethoddefListFields()->Tuple[FieldDescriptor,value]:..."

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
