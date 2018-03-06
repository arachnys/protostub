package gen

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"reflect"

	"github.com/arachnys/protostub"
)

// in spaces
const indentSize = 4

type generator struct {
	depth int
	bw    *bufio.Writer
}

// used to generate mypy stubs from proto data

// writes a mypy type stub to w, generated from the data in p
func Gen(w io.Writer, p *protostub.ProtoData) error {
	gen := &generator{0, bufio.NewWriter(w)}

	for _, i := range p.Types {
		if err := gen.gen(i); err != nil {
			return err
		}
	}

	gen.bw.Flush()

	return nil
}

func (g *generator) indent() error {
	for i := 0; i < g.depth; i++ {
		for j := 0; j < indentSize; j++ {
			if _, err := g.bw.WriteRune(' '); err != nil {
				return err
			}
		}
	}

	return nil
}

func (g *generator) gen(p protostub.ProtoType) error {
	switch t := p.(type) {
	case *protostub.Message:
		return g.genMessage(t)
	case *protostub.Enum:
		return g.genEnum(t)
	case *protostub.Service:
		return g.genService(t)
	}

	return errors.New(fmt.Sprintf("No generator for type %s", reflect.TypeOf(p).Elem().Name()))
}

func (g *generator) genMessage(m *protostub.Message) error {
	return g.genClass(messageToClass(m))
}

func (g *generator) genEnum(e *protostub.Enum) error {
	return g.genClass(enumToClass(e))
}

func (g *generator) genService(s *protostub.Service) error {
	return g.genClass(serviceToClass(s))
}
