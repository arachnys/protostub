package gen

import (
	"fmt"

	"github.com/arachnys/protostub"
)

// This is more or less the same as the types in types.go, however I would like
// to preserve an amount of distinction between the two - in case of future
// diversion
type classData struct {
	name    string
	members []protostub.Member
	types   []protostub.ProtoType
}

func messageToClass(m *protostub.Message) *classData {
	return &classData{
		name:    m.Typename(),
		members: m.Members,
		types:   m.Types,
	}
}

func enumToClass(e *protostub.Enum) *classData {
	return &classData{
		name:    e.Typename(),
		members: e.Members,
	}
}

// generate a mypy/python class
func (g *generator) genClass(c *classData) error {
	err := g.indent()

	if err != nil {
		return err
	}

	_, err = g.bw.WriteString(fmt.Sprintf("class %s:\n", c.name))

	if err != nil {
		return err
	}

	// we're working on members of this class, so indent - ensure to remove
	// the indent when done
	g.depth++

	defer func() {
		g.depth--
	}()

	for n, i := range c.members {
		err := g.indent()

		if err != nil {
			return err
		}

		_, err = g.bw.WriteString(fmt.Sprintf("%s: %s", i.Name(), i.Typename()))

		if err != nil {
			return err
		}

		if n < len(c.members)-1 {
			_, err = g.bw.WriteRune('\n')

			if err != nil {
				return err
			}
		}
	}

	for _, i := range c.types {
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
