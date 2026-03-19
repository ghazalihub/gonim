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
	for i, node := range f.Nodes {
		if i > 0 {
			sb.WriteString("\n")
		}
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
	Pragmas    []string
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
	if p.ReturnTyp != "" && p.ReturnTyp != "void" {
		sb.WriteString(": ")
		sb.WriteString(p.ReturnTyp)
	}
	if len(p.Pragmas) > 0 {
		sb.WriteString(" {.")
		sb.WriteString(strings.Join(p.Pragmas, ", "))
		sb.WriteString(".}")
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
	return strings.TrimSuffix(sb.String(), "\n")
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

type InfixExpr struct {
	Left  string
	Op    string
	Right string
}

func (e *InfixExpr) Render(n int) string {
	return fmt.Sprintf("(%s %s %s)", e.Left, e.Op, e.Right)
}

type PrefixExpr struct {
	Op    string
	Right string
}

func (e *PrefixExpr) Render(n int) string {
	return fmt.Sprintf("%s%s", e.Op, e.Right)
}

type CallExpr struct {
	Fun  string
	Args []string
}

func (c *CallExpr) Render(n int) string {
	return fmt.Sprintf("%s(%s)", c.Fun, strings.Join(c.Args, ", "))
}

type Stmt struct {
	Content string
}

func (s *Stmt) Render(n int) string {
	if strings.TrimSpace(s.Content) == "" {
		return ""
	}
	lines := strings.Split(s.Content, "\n")
	var res []string
	for _, line := range lines {
		if strings.TrimSpace(line) == "" {
			res = append(res, "")
		} else {
			res = append(res, ind(n)+line)
		}
	}
	return strings.Join(res, "\n")
}

type TypeDecl struct {
	Name       string
	TypeParams []string
	Content    string
	Exported   bool
	IsAlias    bool
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

func (i *IfStmt) renderElseOnly(n int) string {
	var sb strings.Builder
	if len(i.Else) == 0 {
		return ""
	}
	if len(i.Else) == 1 {
		if eif, ok := i.Else[0].(*IfStmt); ok {
			sb.WriteString(ind(n))
			sb.WriteString("elif ")
			sb.WriteString(eif.Cond)
			sb.WriteString(":\n")
			for _, node := range eif.Body {
				sb.WriteString(node.Render(n + 1))
				sb.WriteString("\n")
			}
			sb.WriteString(eif.renderElseOnly(n))
			return sb.String()
		}
	}
	sb.WriteString(ind(n))
	sb.WriteString("else:\n")
	for _, node := range i.Else {
		sb.WriteString(node.Render(n + 1))
		sb.WriteString("\n")
	}
	return sb.String()
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
		// Check if it's an 'elif' (the only node in Else is an IfStmt)
		if len(i.Else) == 1 {
			if eif, ok := i.Else[0].(*IfStmt); ok {
				sb.WriteString(ind(n))
				sb.WriteString("elif ")
				sb.WriteString(eif.Cond)
				sb.WriteString(":\n")
				for _, node := range eif.Body {
					sb.WriteString(node.Render(n + 1))
					sb.WriteString("\n")
				}
				if len(eif.Else) > 0 {
					// Recursively render Else of the elif
					nestedIf := &IfStmt{Else: eif.Else}
					sb.WriteString(nestedIf.renderElseOnly(n))
				}
				return strings.TrimSuffix(sb.String(), "\n")
			}
		}
		sb.WriteString(ind(n))
		sb.WriteString("else:\n")
		for _, node := range i.Else {
			sb.WriteString(node.Render(n + 1))
			sb.WriteString("\n")
		}
	}
	return strings.TrimSuffix(sb.String(), "\n")
}

type BlockStmt struct {
	Label string
	Body  []Node
}

func (b *BlockStmt) Render(n int) string {
	var sb strings.Builder
	if b.Label != "" {
		sb.WriteString(ind(n))
		sb.WriteString("block ")
		sb.WriteString(b.Label)
		sb.WriteString(":\n")
		n++
	}
	for _, node := range b.Body {
		sb.WriteString(node.Render(n))
		sb.WriteString("\n")
	}
	return strings.TrimSuffix(sb.String(), "\n")
}

type ForStmt struct {
	Cond string
	Body []Node
}

func (f *ForStmt) Render(n int) string {
	var sb strings.Builder
	sb.WriteString(ind(n))
	if f.Cond == "" {
		sb.WriteString("while true:\n")
	} else {
		sb.WriteString("while ")
		sb.WriteString(f.Cond)
		sb.WriteString(":\n")
	}
	if len(f.Body) == 0 {
		sb.WriteString(ind(n + 1))
		sb.WriteString("discard\n")
	} else {
		for _, node := range f.Body {
			sb.WriteString(node.Render(n + 1))
			sb.WriteString("\n")
		}
	}
	return strings.TrimSuffix(sb.String(), "\n")
}

type CaseStmt struct {
	Expr    string
	Clauses []CaseClause
}

type CaseClause struct {
	Vals []string
	Body []Node
}

func (c *CaseStmt) Render(n int) string {
	var sb strings.Builder
	sb.WriteString(ind(n))
	sb.WriteString("case ")
	sb.WriteString(c.Expr)
	sb.WriteString("\n")
	for _, cl := range c.Clauses {
		sb.WriteString(ind(n))
		if len(cl.Vals) == 0 {
			sb.WriteString("else:\n")
		} else {
			sb.WriteString("of ")
			sb.WriteString(strings.Join(cl.Vals, ", "))
			sb.WriteString(":\n")
		}
		if len(cl.Body) == 0 {
			sb.WriteString(ind(n + 1))
			sb.WriteString("discard\n")
		} else {
			for _, node := range cl.Body {
				sb.WriteString(node.Render(n + 1))
				sb.WriteString("\n")
			}
		}
	}
	return strings.TrimSuffix(sb.String(), "\n")
}
