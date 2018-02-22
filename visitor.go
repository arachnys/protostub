package protostub

import (
	"bufio"
	"fmt"
	"io"
	"strings"

	"github.com/emicklei/proto"
)

type Visitor struct {
	w       io.Writer
	bw      *bufio.Writer
	imports map[string]bool

	// depth is incremented for each class definition
	depth   int
	classes []*class
}

func NewVisitor(w io.Writer) *Visitor {
	return &Visitor{
		w:       w,
		bw:      bufio.NewWriter(w),
		depth:   0,
		imports: make(map[string]bool),
		classes: make([]*class, 0),
	}

}

func (v *Visitor) writeDepth(w io.Writer) {
	for i := 0; i < v.depth; i++ {
		w.Write([]byte("\t"))
	}
}

func (v *Visitor) VisitMessage(m *proto.Message) {
	v.writeDepth(v.bw)

	if len(m.Elements) == 0 {
		v.bw.WriteString(fmt.Sprintf("%s: Any", m.Name))
		return
	}

	if strings.Contains(m.Name, ".") {
		return
	}

	if v.depth == 0 {
		v.bw.WriteString("\n")
	}

	cv := newClassVisitor(v, m.Name)
	v.classes = append(v.classes, cv.class)

	for _, i := range m.Elements {
		i.Accept(cv)
	}

	v.bw.WriteString(cv.String())

	v.bw.Flush()
}

func (v *Visitor) VisitService(s *proto.Service) {
	v.writeDepth(v.bw)

	if v.depth == 0 {
		v.bw.WriteString("\n")
	}

	v.bw.WriteString(fmt.Sprintf("class %sService:\n", s.Name))

	v.depth++

	defer func() {
		v.depth--
	}()

	for _, i := range s.Elements {
		i.Accept(v)
	}

	v.bw.Flush()
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

	v.bw.WriteString(fmt.Sprintf("from %s import *\n", path))
}

func (v *Visitor) VisitNormalField(n *proto.NormalField) {
	v.writeDepth(v.bw)

	if !n.Repeated {
		v.bw.WriteString(fmt.Sprintf("%s: %s\n", n.Name, TranslateType(n.Type)))
	} else {
		v.bw.WriteString(fmt.Sprintf("%s: List[%s]\n", n.Name, TranslateType(n.Type)))
	}
}

func (v *Visitor) VisitEnumField(e *proto.EnumField) {
}

func (v *Visitor) VisitEnum(e *proto.Enum) {
}

func (v *Visitor) VisitComment(c *proto.Comment) {
}

func (v *Visitor) VisitOneof(o *proto.Oneof) {
}

func (v *Visitor) VisitOneofField(o *proto.OneOfField) {
}

func (v *Visitor) VisitReserved(r *proto.Reserved) {
}

func (v *Visitor) VisitRPC(r *proto.RPC) {
	v.writeDepth(v.bw)

	v.bw.WriteString(fmt.Sprintf(
		"def %s(%s) -> %s: ...\n",
		r.Name,
		TranslateType(r.RequestType),
		TranslateType(r.ReturnsType),
	))
}

func (v *Visitor) VisitMapField(m *proto.MapField) {
}

func (v *Visitor) VisitGroup(g *proto.Group) {
}

func (v *Visitor) VisitExtensions(e *proto.Extensions) {
}

var typeMap = map[string]string{
	"int32": "int",
	"int64": "int",
}

func TranslateType(t string) string {
	if val, ok := typeMap[t]; ok {
		return val
	}

	split := strings.Split(t, ".")

	return split[len(split)-1]
}
