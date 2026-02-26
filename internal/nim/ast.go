package nim

import (
	"fmt"
	"strings"
)

type Node interface {
	Render(indent int) string
}

func ind(n int) string {
	return strings.Repeat("  ", n)
}

type File struct {
	Nodes []Node
}

func (f *File) Render(n int) string {
	var sb strings.Builder
	for _, node := range f.Nodes {
		sb.WriteString(node.Render(n))
		sb.WriteString("\n")
	}
	return sb.String()
}

type ProcDecl struct {
	Name       string
	TypeParams []string
	Args       []Arg
	ReturnTyp  string
	Body       []Node
	Exported   bool
}

type Arg struct {
	Name string
	Typ  string
}

func (p *ProcDecl) Render(n int) string {
	var sb strings.Builder
	sb.WriteString(ind(n))
	sb.WriteString("proc ")
	sb.WriteString(p.Name)
	if p.Exported {
		sb.WriteString("*")
	}
	if len(p.TypeParams) > 0 {
		sb.WriteString("[")
		sb.WriteString(strings.Join(p.TypeParams, ", "))
		sb.WriteString("]")
	}
	sb.WriteString("(")
	for i, arg := range p.Args {
		if i > 0 {
			sb.WriteString(", ")
		}
		sb.WriteString(arg.Name)
		sb.WriteString(": ")
		sb.WriteString(arg.Typ)
	}
	sb.WriteString(")")
	if p.ReturnTyp != "" {
		sb.WriteString(": ")
		sb.WriteString(p.ReturnTyp)
	}
	sb.WriteString(" =\n")
	if len(p.Body) == 0 {
		sb.WriteString(ind(n + 1))
		sb.WriteString("discard\n")
	} else {
		for _, node := range p.Body {
			sb.WriteString(node.Render(n + 1))
			sb.WriteString("\n")
		}
	}
	return sb.String()
}

type VarDecl struct {
	Name  string
	Typ   string
	Value string
	Kind  string // "var", "let", "const"
}

func (v *VarDecl) Render(n int) string {
	var sb strings.Builder
	sb.WriteString(ind(n))
	sb.WriteString(v.Kind)
	sb.WriteString(" ")
	sb.WriteString(v.Name)
	if v.Typ != "" {
		sb.WriteString(": ")
		sb.WriteString(v.Typ)
	}
	if v.Value != "" {
		sb.WriteString(" = ")
		sb.WriteString(v.Value)
	}
	return sb.String()
}

type CallExpr struct {
	Fun  string
	Args []string
}

func (c *CallExpr) Render(n int) string {
	return fmt.Sprintf("%s%s(%s)", ind(n), c.Fun, strings.Join(c.Args, ", "))
}

type Stmt struct {
	Content string
}

func (s *Stmt) Render(n int) string {
	return ind(n) + s.Content
}

type TypeDecl struct {
	Name       string
	TypeParams []string
	Content    string
	Exported   bool
}

func (t *TypeDecl) Render(n int) string {
	name := t.Name
	if t.Exported {
		name += "*"
	}
	tp := ""
	if len(t.TypeParams) > 0 {
		tp = "[" + strings.Join(t.TypeParams, ", ") + "]"
	}
	return fmt.Sprintf("%stype %s%s = %s", ind(n), name, tp, t.Content)
}

type ImportStmt struct {
	Pkg string
}

func (i *ImportStmt) Render(n int) string {
	return ind(n) + "import " + i.Pkg
}

type IfStmt struct {
	Cond string
	Body []Node
	Else []Node
}

func (i *IfStmt) Render(n int) string {
	var sb strings.Builder
	sb.WriteString(ind(n))
	sb.WriteString("if ")
	sb.WriteString(i.Cond)
	sb.WriteString(":\n")
	if len(i.Body) == 0 {
		sb.WriteString(ind(n + 1))
		sb.WriteString("discard\n")
	} else {
		for _, node := range i.Body {
			sb.WriteString(node.Render(n + 1))
			sb.WriteString("\n")
		}
	}
	if len(i.Else) > 0 {
		sb.WriteString(ind(n))
		sb.WriteString("else:\n")
		for _, node := range i.Else {
			sb.WriteString(node.Render(n + 1))
			sb.WriteString("\n")
		}
	}
	return strings.TrimSuffix(sb.String(), "\n")
}
