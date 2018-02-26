package protostub

import (
	"io"
	"strings"

	"github.com/emicklei/proto"
)

// This is the main visitor, which passes off to the other ones.
// I try to have a visitor per "major" type - message, service, enum, etc.
// This base Visitor handles the basic stuff.
type Visitor struct {
	// why not set? please?
	imports map[string]bool

	// depth is incremented for each class definition
	depth int
	types []ProtoType
}

func NewVisitor(w io.Writer) *Visitor {
	return &Visitor{
		depth:   0,
		imports: make(map[string]bool),
	}

}

// TODO: make sure everything is here
var typeMap = map[string]string{
	"int32": "int",
	"int64": "int",
}

// takes a protobuf primitive type, makes it a python primitive type
func TranslateType(t string) string {
	if val, ok := typeMap[t]; ok {
		return val
	}

	split := strings.Split(t, ".")

	return split[len(split)-1]
}

func (v *Visitor) writeDepth(w io.Writer) {
	for i := 0; i < v.depth; i++ {
		w.Write([]byte("\t"))
	}
}

func (v *Visitor) VisitMessage(m *proto.Message) {
	mv := NewMessageVisitor(m.Name)

	for _, i := range m.Elements {
		i.Accept(mv)
	}

	v.types = append(v.types, mv.message)
}

func (v *Visitor) VisitSyntax(s *proto.Syntax) {
}

func (v *Visitor) VisitPackage(p *proto.Package) {
}

func (v *Visitor) VisitOption(o *proto.Option) {
}

func (v *Visitor) VisitImport(i *proto.Import) {
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

func (v *Visitor) VisitNormalField(n *proto.NormalField) {
}

func (v *Visitor) VisitEnumField(e *proto.EnumField) {
	panic("Cannot visit enum field")
}

func (v *Visitor) VisitEnum(e *proto.Enum) {
}

func (v *Visitor) VisitComment(c *proto.Comment) {
}

func (v *Visitor) VisitOneof(o *proto.Oneof) {
	panic("Cannot visit oneof")
}

func (v *Visitor) VisitOneofField(o *proto.OneOfField) {
	panic("Cannot visit oneof field")
}

func (v *Visitor) VisitReserved(r *proto.Reserved) {
}

func (v *Visitor) VisitService(r *proto.Service) {
}

func (v *Visitor) VisitRPC(r *proto.RPC) {
}

func (v *Visitor) VisitMapField(m *proto.MapField) {
}

func (v *Visitor) VisitGroup(g *proto.Group) {
}

func (v *Visitor) VisitExtensions(e *proto.Extensions) {
}
