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
	// using the Any type because of this:
	// https://github.com/python/typeshed/blob/master/stdlib/3.4/enum.pyi#L31
	ev.Enum.Members = append(ev.Enum.Members, Member{e.Name, "Any", nil})
}

func (ev *EnumVisitor) VisitEnum(e *proto.Enum) {
	ev.Enum.name = e.Name

	for _, i := range e.Elements {
		i.Accept(ev)
	}
}
