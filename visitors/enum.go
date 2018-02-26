package protostub

import (
	"github.com/emicklei/proto"
)

type EnumVisitor struct {
	Visitor
	message *Message
	members []Member
}

func (ev *EnumVisitor) VisitEnumField(e *proto.EnumField) {
	// all enum members are of type bool in python
	// this is mostly just because they only need to exist, and not store any
	// actual value
	ev.members = append(ev.members, Member{e.Name, "bool"})
}

func (ev *EnumVisitor) VisitEnum(e *proto.Enum) {
	for _, i := range e.Elements {
		i.Accept(ev)
	}
}
