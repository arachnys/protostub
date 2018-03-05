package gen

import (
	"fmt"

	"github.com/arachnys/protostub"
)

// This is more or less the same as the types in types.go, however I would like
// to preserve an amount of distinction between the two - in case of future
// diversion
type classData struct {
	name     string
	members  []protostub.Member
	types    []protostub.ProtoType
	extend   bool
	comments []string
}

func messageToClass(m *protostub.Message) *classData {
	return &classData{
		name:     m.Typename(),
		members:  m.Members,
		types:    m.Types,
		extend:   m.IsExtend,
		comments: m.Comment,
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
	_, err := g.bw.WriteRune('\n')

	if err != nil {
		return err
	}

	if c.extend {
		fmt.Println("Extensions are not yet supported")
		return nil
	}

	if len(c.comments) > 0 {
		for _, i := range c.comments {
			g.indent()
			g.bw.WriteString(fmt.Sprintf("#%s\n", i))
		}
	}

	err = g.indent()

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
		for _, j := range i.Comment {
			g.indent()
			g.bw.WriteString(fmt.Sprintf("#%s\n", j))
		}

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

	// let's make that constructor
	g.bw.WriteRune('\n')
	g.indent()
	g.bw.WriteString("def __init__(self, ")

	for n, i := range c.members {
		if n < len(c.members)-1 {
			g.bw.WriteString(fmt.Sprintf("%s: %s = None, ", i.Name(), i.Typename()))
			continue
		}

		g.bw.WriteString(fmt.Sprintf("%s: %s = None) -> None: ...\n", i.Name(), i.Typename()))
	}

	for _, i := range c.types {
		// enums need to be treated differently
		if e, ok := i.(*protostub.Enum); ok {
			for _, j := range e.Members {
				g.indent()
				g.bw.WriteString(fmt.Sprintf("%s: Any\n", j.Name()))
			}
		}

		err := g.indent()

		if err != nil {
			return err
		}

		err = g.gen(i)

		if err != nil {
			return err
		}
	}

	// then we just need to generate all the default methods that protoc adds to
	// python classes
	defaults := []string{fmt.Sprintf("CopyFrom(self, other: %s) -> Any", c.name), "ListFields() -> Tuple[FieldDescriptor, value]"}

	for _, i := range defaults {
		g.indent()
		g.bw.WriteString(fmt.Sprintf("def %s: ...\n", i))
	}

	return nil
}
