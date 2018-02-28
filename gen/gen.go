package gen

import (
	"bufio"
	"io"

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
	}

	return nil
}

func (g *generator) genMessage(m *protostub.Message) error {
	return nil
}
