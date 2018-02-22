package protostub

import (
	"bytes"
	"fmt"

	"github.com/emicklei/proto"
)

type member struct {
	name     string
	typename string
}

type class struct {
	name    string
	members []member
}

// As soon as the main visitor finds a messages, it creates an instance of this.
// It is mostly the same, however keeps track of the members and subclasses it
// finds. This is so that it can generate a constructor at the end, as well
// as any other methods. It's also nice to keep track of what types there are,
// in at least some way.
type classVisitor struct {
	*Visitor
	class   *class
	classes []*classVisitor
}

func newClassVisitor(v *Visitor, name string) *classVisitor {
	cv := &classVisitor{
		Visitor: v,
		class: &class{
			name:    name,
			members: make([]member, 0),
		},
		classes: make([]*classVisitor, 0),
	}

	return cv
}

func (cv *classVisitor) addMember(m member) {
	cv.class.members = append(cv.class.members, m)
}

func (cv *classVisitor) VisitNormalField(n *proto.NormalField) {
	name := n.Name
	var typename string

	if !n.Repeated {
		typename = TranslateType(n.Type)
	} else {
		typename = fmt.Sprintf("List[%s]", TranslateType(n.Type))
	}

	cv.addMember(member{
		name,
		typename,
	})
}

func (cv *classVisitor) VisitMessage(m *proto.Message) {

	scv := newClassVisitor(cv.Visitor, m.Name)

	for _, i := range m.Elements {
		i.Accept(scv)
	}

	cv.classes = append(cv.classes, scv)
}

func (cv *classVisitor) String() string {
	b := bytes.NewBuffer(nil)

	b.WriteString(fmt.Sprintf("class %s:\n", cv.class.name))

	cv.Visitor.depth++

	defer func() {
		cv.Visitor.depth--
	}()

	for _, i := range cv.classes {
		cv.Visitor.writeDepth(b)
		b.WriteString(i.String())
	}

	for _, i := range cv.class.members {
		cv.Visitor.writeDepth(b)
		b.WriteString(fmt.Sprintf("%s: %s\n", i.name, i.typename))
	}

	// generate the constructor
	b.WriteRune('\n')
	cv.Visitor.writeDepth(b)
	b.WriteString("def __init__(self, ")

	for n, i := range cv.class.members {
		if n < len(cv.class.members)-1 {
			b.WriteString(fmt.Sprintf("%s: %s = None, ", i.name, i.typename))
			continue
		}

		b.WriteString(fmt.Sprintf("%s: %s = None) -> None: ...\n", i.name, i.typename))
	}

	// add the methods protobuf adds
	cv.writeDepth(b)
	b.WriteString(fmt.Sprintf("def CopyFrom(self, other: %s) -> None: ...\n", cv.class.name))
	cv.writeDepth(b)
	b.WriteString("def ListFields(self) -> Tuple[FieldDescriptor, value]: ...\n")

	// just to make it all a bit more readable
	b.WriteRune('\n')

	return b.String()
}

func (cv *classVisitor) VisitOneof(o *proto.Oneof) {
	for _, i := range o.Elements {
		i.Accept(cv)
	}
}

func (cv *classVisitor) VisitOneofField(o *proto.OneOfField) {
	cv.VisitNormalField(&proto.NormalField{
		Field: o.Field,
	})
}

func (cv *classVisitor) VisitEnumField(e *proto.EnumField) {
	cv.addMember(member{
		e.Name,
		"bool",
	})
}

func (v *classVisitor) VisitEnum(e *proto.Enum) {
	for _, i := range e.Elements {
		i.Accept(v)
	}
}
