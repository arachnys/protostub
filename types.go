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
}

func (m Member) Name() string     { return m.name }
func (m Member) Typename() string { return m.typename }

func (m Message) Name() string     { return m.name }
func (m Message) Typename() string { return m.name }

func (s Service) Name() string     { return s.name }
func (s Service) Typename() string { return s.name }

// TODO: Move this to be part of the gen package
//func (cv *classVisitor) String() string {
//	b := bytes.NewBuffer(nil)
//
//	b.WriteString(fmt.Sprintf("class %s:\n", cv.class.name))
//
//	cv.Visitor.depth++

//	defer func() {
//		cv.Visitor.depth--
//	}()

//	for _, i := range cv.classes {
//		cv.Visitor.writeDepth(b)
//		b.WriteString(i.String())
//	}

//	for _, i := range cv.class.members {
//		cv.Visitor.writeDepth(b)
//		b.WriteString(fmt.Sprintf("%s: %s\n", i.name, i.typename))
//	}

//	// generate the constructor
//	b.WriteRune('\n')
//	cv.Visitor.writeDepth(b)
//	b.WriteString("def __init__(self, ")

//	for n, i := range cv.class.members {
//		if n < len(cv.class.members)-1 {
//			b.WriteString(fmt.Sprintf("%s: %s = None, ", i.name, i.typename))
//			continue
//		}

//		b.WriteString(fmt.Sprintf("%s: %s = None) -> None: ...\n", i.name, i.typename))
//	}

// add the methods protobuf adds
//	cv.writeDepth(b)
//	b.WriteString(fmt.Sprintf("def CopyFrom(self, other: %s) -> None: ...\n", cv.class.name))
//	cv.writeDepth(b)
//	b.WriteString("def ListFields(self) -> Tuple[FieldDescriptor, value]: ...\n")

// just to make it all a bit more readable
//	b.WriteRune('\n')

//	return b.String()
//}
