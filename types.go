package protostub

import (
	"fmt"
	"strings"
)

// This contains the types that the in-memory representation of the protobuf
// will be built out of. Some have associated visitors that will generate them,
// while others are just components

// ProtoType defines that a type must have both a name and a TypeName
type ProtoType interface {
	Name() string
	Typename() string
}

// Member is a member of a message
type Member struct {
	name     string
	typename string
	Comment  []string
}

// Function is a function, usually a RPC method in this case
type Function struct {
	name       string
	returnType string
	parameters []string
	Comment    []string
}

// Message represents a ProtoBuf message
type Message struct {
	name     string
	Types    []ProtoType
	Members  []Member
	IsExtend bool
	Comment  []string
}

// Service representts a protobuf service
type Service struct {
	name      string
	Types     []ProtoType
	Functions []Function
	Comment   []string
}

// Enum is a protobuf enum
type Enum struct {
	name    string
	Members []Member
	Values  []int
}

// Name returns the name of the member
func (m Member) Name() string { return m.name }

// Typename returns the type name of the member
func (m Member) Typename() string { return m.typename }

// Name returns the name of the message
func (m Message) Name() string { return m.name }

// Typename returns the type name of the message. Note that in this case,
// Name() == Typename()
func (m Message) Typename() string { return m.name }

// Name returns the name of the service
func (s Service) Name() string { return s.name }

// Typename returns the typename of the service. Like Message, in this case
// Typename is the same as Name.
func (s Service) Typename() string { return s.name }

// Name returns the name of the enum.
func (s Enum) Name() string { return s.name }

// Typename returns the typename of the enum. Again, it is the same as Name.
func (s Enum) Typename() string { return s.name }

// Name returns the name of the function.
func (f Function) Name() string { return f.name }

// Typename returns the typename of the function. This includes the name,
// parameter, and return value in the form foo(bar) -> baz.
func (f Function) Typename() string {
	return fmt.Sprintf("%s(%s) -> %s", f.name, strings.Join(f.parameters, ", "), f.returnType)
}
