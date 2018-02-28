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
		err := gen.gen(i)

		if err != nil {
			return err
		}
	}

	gen.bw.Flush()

	return nil
}

func (g *generator) indent() error {
	for i := 0; i < g.depth; i++ {
		for j := 0; j < indentSize; j++ {
			_, err := g.bw.WriteRune(' ')

			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (g *generator) gen(p protostub.ProtoType) error {
	switch t := p.(type) {
	case *protostub.Message:
		err := g.genMessage(t)

		if err != nil {
			return err
		}

		return nil
	}

	return errors.New(fmt.Sprintf("No generator for type %s", reflect.TypeOf(p).Elem().Name()))
}

func (g *generator) genMessage(m *protostub.Message) error {
	err := g.indent()

	if err != nil {
		return err
	}

	_, err = g.bw.WriteString(fmt.Sprintf("class %s:\n", m.Name()))

	if err != nil {
		return err
	}

	// we're working on members of this message, so indent - ensure to remove
	// the indent when done
	g.depth++

	defer func() {
		g.depth--
	}()

	for n, i := range m.Members {
		err := g.indent()

		if err != nil {
			return err
		}

		_, err = g.bw.WriteString(fmt.Sprintf("%s: %s", i.Name(), i.Typename()))

		if err != nil {
			return err
		}

		if n < len(m.Members)-1 {
			_, err = g.bw.WriteRune('\n')

			if err != nil {
				return err
			}
		}
	}

	for _, i := range m.Types {
		err := g.indent()

		if err != nil {
			return err
		}

		err = g.gen(i)

		if err != nil {
			return err
		}
	}

	return nil
}
