package protostub

import (
	"github.com/emicklei/proto"
)

// EnumVisitor is the Visitor for an enum, created by the main visitor. It
// contains an Enum with all the type data - see types.go.
type EnumVisitor struct {
	ProtoData
	Enum Enum
}

// NewEnumVisitor creates a new EnumVisitor
func NewEnumVisitor() *EnumVisitor {
	return &EnumVisitor{
		ProtoData: ProtoData{},
	}
}

// VisitEnumField creates an enum field to the member list
func (ev *EnumVisitor) VisitEnumField(e *proto.EnumField) {
	// protobufs only uses integers for enum values
	ev.Enum.Members = append(ev.Enum.Members, Member{e.Name, ev.Enum.name, nil})
	ev.Enum.Values = append(ev.Enum.Values, e.Integer)
}

// VisitEnum will set the correct name and visit the elements
func (ev *EnumVisitor) VisitEnum(e *proto.Enum) {
	ev.Enum.name = e.Name

	for _, i := range e.Elements {
		i.Accept(ev)
	}
}
