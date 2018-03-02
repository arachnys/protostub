package protostub

import (
	"github.com/emicklei/proto"
)

type EnumVisitor struct {
	ProtoData
	Enum Enum
}

func NewEnumVisitor() *EnumVisitor {
	return &EnumVisitor{
		ProtoData: ProtoData{},
	}
}

func (ev *EnumVisitor) VisitEnumField(e *proto.EnumField) {
	// all enum members are of type bool in python
	// this is mostly just because they only need to exist, and not store any
	// actual value
	ev.Enum.Members = append(ev.Enum.Members, Member{e.Name, "bool"})
}

func (ev *EnumVisitor) VisitEnum(e *proto.Enum) {
	ev.Enum.name = e.Name

	for _, i := range e.Elements {
		i.Accept(ev)
	}
}
