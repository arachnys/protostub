package protostub

import (
	"fmt"

	"github.com/emicklei/proto"
)

type MessageVisitor struct {
	ProtoData
	message *Message
}

func NewMessageVisitor(name string, extend bool, comment *proto.Comment) *MessageVisitor {

	cv := &MessageVisitor{
		ProtoData: ProtoData{},
		message: &Message{
			name:     name,
			Types:    make([]ProtoType, 0),
			IsExtend: extend,
			Comment:  comment.Lines,
		},
	}

	return cv
}

func (mv *MessageVisitor) addMember(m Member) {
	mv.message.Members = append(mv.message.Members, m)
}

func (mv *MessageVisitor) VisitNormalField(n *proto.NormalField) {
	name := n.Name
	var typename string

	// repeated = it's a list
	if !n.Repeated {
		typename = TranslateType(n.Type)
	} else {
		typename = fmt.Sprintf("List[%s]", TranslateType(n.Type))
	}

	if n.Comment == nil {
		n.Comment = &proto.Comment{Lines: make([]string, 0)}
	}

	mv.addMember(Member{
		name:     name,
		typename: typename,
		Comment:  n.Comment.Lines,
	})

}

func (mv *MessageVisitor) VisitOneof(o *proto.Oneof) {
	for _, i := range o.Elements {
		i.Accept(mv)
	}
}

// look into some sort of variant instead
func (mv *MessageVisitor) VisitOneofField(o *proto.OneOfField) {
	mv.VisitNormalField(&proto.NormalField{
		Field: o.Field,
	})
}
