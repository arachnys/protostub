package protostub

import (
	"io"
	"strings"

	"github.com/emicklei/proto"
)

// ProtoData is the main visitor, which passes off to the other ones.
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

// New creates a new protodata object. It takes a reader to the proto file
func New(r io.Reader) *ProtoData {
	return &ProtoData{
		Types:   make([]ProtoType, 0),
		depth:   0,
		imports: make(map[string]bool),
		r:       r,
	}
}

// Parse will parse the proto file, and dispatch to all the visitors. Error is
// returned on failure
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

// TranslateType takes in a protobuf type and translates it to a Python primitive
// type. If there is no translation needed, the input is returned unaltered.
func TranslateType(t string) string {
	if val, ok := typeMap[t]; ok {
		return val
	}

	split := strings.Split(t, ".")

	return split[len(split)-1]
}

// VisitMessage will create a MessageVisitor, and dispatch it. The message type
// is included in protodata.
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

// VisitSyntax presently does nothing.
func (v *ProtoData) VisitSyntax(s *proto.Syntax) {
}

// VisitPackage presently does nothing.
func (v *ProtoData) VisitPackage(p *proto.Package) {
}

// VisitOption presently does nothing.
func (v *ProtoData) VisitOption(o *proto.Option) {
}

// VisitImport will try to translate a proto import to a python one. Currently it
// is a bit of a hack.
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

// VisitNormalField does nothing here, and is implemented in the message visitor.
func (v *ProtoData) VisitNormalField(n *proto.NormalField) {
}

// VisitEnumField will panic when called here, as enum fields should be inside
// enums.
func (v *ProtoData) VisitEnumField(e *proto.EnumField) {
	panic("Cannot visit enum field")
}

// VisitEnum will create a new Enum visitor and dispatch it.
func (v *ProtoData) VisitEnum(e *proto.Enum) {
	ev := NewEnumVisitor()

	e.Accept(ev)

	v.Types = append(v.Types, &ev.Enum)
}

// VisitComment currently does nothing, comments are only handled when attached
// to messages and their fields.
func (v *ProtoData) VisitComment(c *proto.Comment) {
}

// VisitOneof will panic when called from here.
func (v *ProtoData) VisitOneof(o *proto.Oneof) {
	panic("Cannot visit oneof")
}

// VisitOneofField will panic when called from here.
func (v *ProtoData) VisitOneofField(o *proto.OneOfField) {
	panic("Cannot visit oneof field")
}

// VisitReserved is currently not implemented.
func (v *ProtoData) VisitReserved(r *proto.Reserved) {
}

// VisitService will create a new service visitor, and dispatch it.
func (v *ProtoData) VisitService(r *proto.Service) {
	if r.Comment == nil {
		r.Comment = &proto.Comment{Lines: make([]string, 0)}
	}

	sv := NewServiceVisitor(r.Name, r.Comment)

	for _, i := range r.Elements {
		i.Accept(sv)
	}

	v.Types = append(v.Types, sv.service)
	sv.service.Types = sv.Types
}

// VisitRPC will currently do nothing, as RPC should be inside a service.
func (v *ProtoData) VisitRPC(r *proto.RPC) {
}

// VisitMapField is currently not implemented
func (v *ProtoData) VisitMapField(m *proto.MapField) {
}

// VisitGroup is currently not implemented
func (v *ProtoData) VisitGroup(g *proto.Group) {
}

// VisitExtensions is currently not implemented
func (v *ProtoData) VisitExtensions(e *proto.Extensions) {
}
