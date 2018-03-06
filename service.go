package protostub

import (
	"github.com/emicklei/proto"
)

// ServiceVisitor will visit a service and extract the types.
type ServiceVisitor struct {
	ProtoData
	service *Service
}

// NewServiceVisitor will create a new service visitor, from a name and a comment.
func NewServiceVisitor(name string, comment *proto.Comment) *ServiceVisitor {
	sv := &ServiceVisitor{
		ProtoData: ProtoData{},
		service: &Service{
			name:      name,
			Types:     make([]ProtoType, 0),
			Functions: make([]Function, 0),
			Comment:   comment.Lines,
		},
	}

	return sv
}

func (sv *ServiceVisitor) addFunction(f Function) {
	sv.service.Functions = append(sv.service.Functions, f)
}

// VisitRPC will visit a RPC, and make sure the function specification is included
// in the service data definition.
func (sv *ServiceVisitor) VisitRPC(r *proto.RPC) {
	if r.Comment == nil {
		r.Comment = &proto.Comment{
			Lines: make([]string, 0),
		}
	}

	f := Function{
		name:       r.Name,
		returnType: r.ReturnsType,
		parameters: []string{r.RequestType},
		Comment:    r.Comment.Lines,
	}

	sv.addFunction(f)
}
