package translator

import (
	"fmt"
	"strings"

	"github.com/user/go2nim/internal/builtins"
	"github.com/user/go2nim/internal/ir"
	"github.com/user/go2nim/internal/nim"
	"github.com/user/go2nim/internal/types"
)

type translator struct {
	imports map[string]bool
}

func (t *translator) AddImport(name string) {
	t.imports[name] = true
}

func Translate(prog *ir.Program) []*nim.File {
	var files []*nim.File
	for _, pkg := range prog.Packages {
		t := &translator{imports: make(map[string]bool)}
		files = append(files, t.translatePackage(pkg))
	}
	return files
}

func (t *translator) translatePackage(pkg *ir.Package) *nim.File {
	nimFile := &nim.File{}
	var nodes []nim.Node
	for _, file := range pkg.Files {
		for _, decl := range file.Decls {
			nimDecl := t.translateDecl(decl)
			if nimDecl != nil {
				nodes = append(nodes, nimDecl)
			}
		}
	}
	for imp := range t.imports {
		nimFile.Nodes = append(nimFile.Nodes, &nim.ImportStmt{Pkg: imp})
	}
	nimFile.Nodes = append(nimFile.Nodes, nodes...)
	return nimFile
}

func (t *translator) translateDecl(decl ir.Decl) nim.Node {
	switch d := decl.(type) {
	case *ir.FuncDecl:
		args := t.translateFields(d.Type.Params)
		if d.Receiver != nil {
			recvName := "self"
			if len(d.Receiver.Names) > 0 {
				recvName = d.Receiver.Names[0]
			}
			args = append([]nim.Arg{{Name: recvName, Typ: types.MapType(d.Receiver.Type, t)}}, args...)
		}
		return &nim.ProcDecl{
			Name:       d.Name,
			TypeParams: t.translateTypeParams(d.Type.TypeParams),
			Args:       args,
			ReturnTyp:  t.translateResults(d.Type.Results),
			Body:       t.translateBlock(d.Body),
			Exported:   isExported(d.Name),
		}
	case *ir.VarDecl:
		kind := "var"
		value := ""
		if len(d.Values) > 0 {
			value = t.translateExpr(d.Values[0])
		}
		if len(d.Embeds) > 0 {
			kind = "const"
			value = fmt.Sprintf("staticRead(\"%s\")", d.Embeds[0])
		}
		return &nim.VarDecl{
			Kind:  kind,
			Name:  d.Names[0], // Simplified
			Typ:   types.MapType(d.Type, t),
			Value: value,
		}
	case *ir.TypeDecl:
		content := ""
		underlying := d.Type
		if nt, ok := d.Type.(*ir.NamedType); ok {
			underlying = nt.Underlying
		}
		if st, ok := underlying.(*ir.StructType); ok {
			var sb strings.Builder
			sb.WriteString("object\n")
			for _, f := range st.Fields {
				for _, name := range f.Names {
					exported := ""
					if isExported(name) {
						exported = "*"
					}
					sb.WriteString(fmt.Sprintf("    %s%s: %s\n", name, exported, types.MapType(f.Type, t)))
				}
			}
			content = sb.String()
		} else if it, ok := underlying.(*ir.InterfaceType); ok {
			var sb strings.Builder
			sb.WriteString("concept x\n")
			for _, m := range it.Methods {
				sb.WriteString(fmt.Sprintf("    x.%s() is %s\n", m.Name, t.translateResults(m.Type.Results)))
			}
			content = sb.String()
		} else {
			content = types.MapType(d.Type, t)
		}
		return &nim.TypeDecl{
			Name:       d.Name,
			TypeParams: t.translateTypeParams(d.TypeParams),
			Content:    content,
			Exported:   isExported(d.Name),
		}
	case *ir.ConstDecl:
		return &nim.VarDecl{
			Kind:  "const",
			Name:  d.Names[0], // Simplified
			Typ:   types.MapType(d.Type, t),
			Value: t.translateExpr(d.Values[0]), // Simplified
		}
	case *ir.CGoDecl:
		return &nim.Stmt{Content: fmt.Sprintf("{.emit: \"\"\"\n%s\n\"\"\".}", d.Preamble)}
	case *ir.ImportDecl:
		if d.Path == "C" {
			return nil
		}
		// Map Go import to Nim.
		// For now just use the last element of the path.
		parts := strings.Split(d.Path, "/")
		pkg := parts[len(parts)-1]
		return &nim.ImportStmt{Pkg: pkg}
	}
	return nil
}

func (t *translator) translateTypeParams(tps []*ir.TypeParam) []string {
	var res []string
	for _, tp := range tps {
		// Nim generics often don't need explicit constraints if they are concepts
		// But we can add them if needed.
		res = append(res, tp.Name)
	}
	return res
}

func (t *translator) translateFields(fields []*ir.Field) []nim.Arg {
	var args []nim.Arg
	for _, f := range fields {
		for _, name := range f.Names {
			args = append(args, nim.Arg{
				Name: name,
				Typ:  types.MapType(f.Type, t),
			})
		}
	}
	return args
}

func (t *translator) translateResults(results []*ir.Field) string {
	if len(results) == 0 {
		return ""
	}
	if len(results) == 1 {
		return types.MapType(results[0].Type, t)
	}
	// Multiple returns in Go -> Tuple in Nim
	var res []string
	for _, r := range results {
		res = append(res, types.MapType(r.Type, t))
	}
	return fmt.Sprintf("(%s)", strings.Join(res, ", "))
}

func (t *translator) translateBlock(b *ir.BlockStmt) []nim.Node {
	if b == nil {
		return nil
	}
	var nodes []nim.Node
	for _, stmt := range b.List {
		nodes = append(nodes, t.translateStmt(stmt))
	}
	return nodes
}

func (t *translator) translateStmt(stmt ir.Stmt) nim.Node {
	switch s := stmt.(type) {
	case *ir.ExprStmt:
		return &nim.Stmt{Content: t.translateExpr(s.X)}
	case *ir.ReturnStmt:
		if len(s.Results) == 0 {
			return &nim.Stmt{Content: "return"}
		}
		var res []string
		for _, r := range s.Results {
			res = append(res, t.translateExpr(r))
		}
		return &nim.Stmt{Content: "return " + strings.Join(res, ", ")}
	case *ir.DeclStmt:
		var nodes []nim.Node
		for _, decl := range s.Decls {
			nodes = append(nodes, t.translateDecl(decl))
		}
		// If multiple, wrap in a block or just return the first for now if Nim doesn't support multiple stmts in one Node.
		// Actually Stmt content can be multiline or we can change translateStmt to return []nim.Node.
		if len(nodes) == 1 {
			return nodes[0]
		}
		var sb strings.Builder
		for _, n := range nodes {
			sb.WriteString(n.Render(0))
			sb.WriteString("\n")
		}
		return &nim.Stmt{Content: strings.TrimSuffix(sb.String(), "\n")}
	case *ir.AssignStmt:
		var lhs, rhs []string
		for _, l := range s.Lhs {
			lhs = append(lhs, t.translateExpr(l))
		}
		for _, r := range s.Rhs {
			rhs = append(rhs, t.translateExpr(r))
		}
		lhsStr := strings.Join(lhs, ", ")
		rhsStr := strings.Join(rhs, ", ")
		if len(lhs) > 1 {
			lhsStr = "(" + lhsStr + ")"
		}
		if len(rhs) > 1 {
			rhsStr = "(" + rhsStr + ")"
		}
		op := s.Op
		if op == ":=" {
			typ := types.MapType(s.Lhs[0].GetType(), t)
			return &nim.Stmt{Content: "var " + lhsStr + ": " + typ + " = " + rhsStr}
		}
		return &nim.Stmt{Content: fmt.Sprintf("%s %s %s", lhsStr, op, rhsStr)}
	case *ir.IfStmt:
		var elseNodes []nim.Node
		if s.Else != nil {
			switch e := s.Else.(type) {
			case *ir.BlockStmt:
				elseNodes = t.translateBlock(e)
			default:
				elseNodes = append(elseNodes, t.translateStmt(e))
			}
		}
		return &nim.IfStmt{
			Cond: t.translateExpr(s.Cond),
			Body: t.translateBlock(s.Body),
			Else: elseNodes,
		}
	case *ir.DeferStmt:
		return &nim.Stmt{Content: "defer: " + t.translateExpr(s.Call)}
	}
	return &nim.Stmt{Content: "discard # unknown stmt"}
}

func (t *translator) translateExpr(expr ir.Expr) string {
	if expr == nil {
		return ""
	}
	switch e := expr.(type) {
	case *ir.Ident:
		return e.Name
	case *ir.BasicLit:
		return e.Value
	case *ir.CallExpr:
		fun := t.translateExpr(e.Fun)
		if fun == "make" && len(e.Args) >= 1 {
			if te, ok := e.Args[0].(*ir.TypeExpr); ok {
				switch te.Type.(type) {
				case *ir.MapType:
					return "initTable[" + types.MapType(te.Type.(*ir.MapType).Key, t) + ", " + types.MapType(te.Type.(*ir.MapType).Value, t) + "]()"
				case *ir.SliceType:
					size := "0"
					if len(e.Args) >= 2 {
						size = t.translateExpr(e.Args[1])
					}
					return "newSeq[" + types.MapType(te.Type.(*ir.SliceType).Elem, t) + "](" + size + ")"
				}
			}
		}

		isMember := false
		if b, ok := builtins.Builtins[fun]; ok {
			fun = b.NimName
			isMember = b.IsMember
		}
		var args []string
		for _, arg := range e.Args {
			args = append(args, t.translateExpr(arg))
		}
		if isMember && len(args) > 0 {
			return fmt.Sprintf("%s.%s(%s)", args[0], fun, strings.Join(args[1:], ", "))
		}
		return fmt.Sprintf("%s(%s)", fun, strings.Join(args, ", "))
	case *ir.BinaryExpr:
		op := e.Op
		// Map Go operators to Nim if they differ
		switch op {
		case "&&":
			op = "and"
		case "||":
			op = "or"
		case "!":
			op = "not" // though this is Unary
		}
		return fmt.Sprintf("(%s %s %s)", t.translateExpr(e.X), op, t.translateExpr(e.Y))
	case *ir.SelectorExpr:
		if x, ok := e.X.(*ir.Ident); ok && x.Name == "C" {
			return e.Sel
		}
		return fmt.Sprintf("%s.%s", t.translateExpr(e.X), e.Sel)
	case *ir.IndexExpr:
		return fmt.Sprintf("%s[%s]", t.translateExpr(e.X), t.translateExpr(e.Index))
	case *ir.UnaryExpr:
		op := e.Op
		if op == "!" {
			op = "not "
		}
		return fmt.Sprintf("%s%s", op, t.translateExpr(e.X))
	case *ir.CompositeLit:
		var elms []string
		for _, elm := range e.Elms {
			elms = append(elms, t.translateExpr(elm))
		}
		typ := types.MapType(e.Type, t)
		return fmt.Sprintf("%s(%s)", typ, strings.Join(elms, ", "))
	case *ir.TypeExpr:
		return types.MapType(e.Type, t)
	}
	return "unknown_expr"
}

func isExported(name string) bool {
	if name == "" {
		return false
	}
	return name[0] >= 'A' && name[0] <= 'Z'
}
