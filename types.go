package protostub

type ProtoType interface {
	Name() string
	Typename() string
}

type Member struct {
	name     string
	typename string
}

type Message struct {
	name    string
	Types   []ProtoType
	Members []Member
}

type Service struct {
	name  string
	types []ProtoType
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
