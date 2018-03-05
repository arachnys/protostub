package protostub

import (
	"fmt"
	"strings"
)

type ProtoType interface {
	Name() string
	Typename() string
}

type Member struct {
	name     string
	typename string
	Comment  []string
}

type Function struct {
	name       string
	returnType string
	parameters []string
	Comment    []string
}

type Message struct {
	name     string
	Types    []ProtoType
	Members  []Member
	IsExtend bool
	Comment  []string
}

type Service struct {
	name      string
	Types     []ProtoType
	Functions []Function
	Comment   []string
}

type Enum struct {
	name    string
	Members []Member
}

func (m Member) Name() string     { return m.name }
func (m Member) Typename() string { return m.typename }

func (m Message) Name() string     { return m.name }
func (m Message) Typename() string { return m.name }

func (s Service) Name() string     { return s.name }
func (s Service) Typename() string { return s.name }

func (s Enum) Name() string     { return s.name }
func (s Enum) Typename() string { return s.name }

func (f Function) Name() string { return f.name }
func (f Function) Typename() string {
	return fmt.Sprintf("%s(%s) -> %s", f.name, strings.Join(f.parameters, ", "), f.returnType)
}
