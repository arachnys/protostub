package protostub

import (
	"io"
	"strings"

	"github.com/emicklei/proto"
)

// This is the main visitor, which passes off to the other ones.
// I try to have a visitor per "major" type - message, service, enum, etc.
// This base Visitor handles the basic stuff.
type ProtoData struct {
	Types []ProtoType

	// why not set? please?
	imports map[string]bool

	// depth is incremented for each class definition
	r     io.Reader
	depth int
}

func New(r io.Reader) *ProtoData {
	return &ProtoData{
		Types:   make([]ProtoType, 0),
		depth:   0,
		imports: make(map[string]bool),
		r:       r,
	}
}

func (v *ProtoData) Parse() error {
	p := proto.NewParser(v.r)

	pr, err := p.Parse()

	if err != nil {
		return err
	}

	for _, i := range pr.Elements {
		i.Accept(v)
	}

	return nil
}

// TODO: make sure everything is here
var typeMap = map[string]string{
	"int32":    "int",
	"int64":    "int",
	"double":   "float",
	"uint32":   "int",
	"uint64":   "int",
	"sint32":   "int",
	"sint64":   "int",
	"fixed32":  "int",
	"fixed64":  "int",
	"sfixed32": "int",
	"sfixed64": "int",
	"bytes":    "str",
}

// takes a protobuf primitive type, makes it a python primitive type
func TranslateType(t string) string {
	if val, ok := typeMap[t]; ok {
		return val
	}

	split := strings.Split(t, ".")

	return split[len(split)-1]
}

func (v *ProtoData) writeDepth(w io.Writer) {
	for i := 0; i < v.depth; i++ {
		w.Write([]byte("\t"))
	}
}

func (v *ProtoData) VisitMessage(m *proto.Message) {
	if m.Comment == nil {
		m.Comment = &proto.Comment{Lines: make([]string, 0)}
	}

	mv := NewMessageVisitor(m.Name, m.IsExtend, m.Comment)

	for _, i := range m.Elements {
		i.Accept(mv)
	}

	v.Types = append(v.Types, mv.message)
	mv.message.Types = mv.Types
}

func (v *ProtoData) VisitSyntax(s *proto.Syntax) {
}

func (v *ProtoData) VisitPackage(p *proto.Package) {
}

func (v *ProtoData) VisitOption(o *proto.Option) {
}

func (v *ProtoData) VisitImport(i *proto.Import) {
	// a bit of a hack, yet perhaps it will work?
	// then remove the specific file, so we just have the path of the package
	split := strings.Split(i.Filename, "/")

	// the path is all but the last
	path := strings.Join(split[:len(split)-1], ".")

	// because we can't figure out the class we are importing by the import
	// path, import it all.

	// if it already exists, don't bother continuing.
	if _, ok := v.imports[path]; ok {
		return
	}

	v.imports[path] = true
}

func (v *ProtoData) VisitNormalField(n *proto.NormalField) {
}

func (v *ProtoData) VisitEnumField(e *proto.EnumField) {
	panic("Cannot visit enum field")
}

func (v *ProtoData) VisitEnum(e *proto.Enum) {
	ev := NewEnumVisitor()

	e.Accept(ev)

	v.Types = append(v.Types, &ev.Enum)
}

func (v *ProtoData) VisitComment(c *proto.Comment) {
}

func (v *ProtoData) VisitOneof(o *proto.Oneof) {
	panic("Cannot visit oneof")
}

func (v *ProtoData) VisitOneofField(o *proto.OneOfField) {
	panic("Cannot visit oneof field")
}

func (v *ProtoData) VisitReserved(r *proto.Reserved) {
}

func (v *ProtoData) VisitService(r *proto.Service) {
}

func (v *ProtoData) VisitRPC(r *proto.RPC) {
}

func (v *ProtoData) VisitMapField(m *proto.MapField) {
}

func (v *ProtoData) VisitGroup(g *proto.Group) {
}

func (v *ProtoData) VisitExtensions(e *proto.Extensions) {
}
