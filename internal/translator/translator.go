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
	t.AddImport("gostdnim/builtin")
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
	if decl == nil {
		return nil
	}
	switch d := decl.(type) {
	case *ir.FuncDecl:
		args := t.translateFields(d.Type.Params, d.Type.Variadic)
		if d.Receiver != nil {
			recvName := "self"
			if len(d.Receiver.Names) > 0 {
				recvName = d.Receiver.Names[0]
			}
			args = append([]nim.Arg{{Name: recvName, Typ: types.MapType(d.Receiver.Type, t)}}, args...)
		}
		return &nim.ProcDecl{
			Name:       EscapeNimKeyword(d.Name),
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
			Name:  EscapeNimKeyword(d.Names[0]), // Simplified
			Typ:   types.MapType(d.Type, t),
			Value: value,
		}
	case *ir.TypeDecl:
		content := ""
		underlying := d.Type
		if nt, ok := d.Type.(*ir.NamedType); ok {
			underlying = nt.Underlying
		}
		if d.Alias {
			content = types.MapType(underlying, t)
		} else if st, ok := underlying.(*ir.StructType); ok {
			var sb strings.Builder
			sb.WriteString("object\n")
			for _, f := range st.Fields {
				for _, name := range f.Names {
					exported := ""
					if isExported(name) {
						exported = "*"
					}
					sb.WriteString(fmt.Sprintf("    %s%s: %s\n", EscapeNimKeyword(name), exported, types.MapType(f.Type, t)))
				}
			}
			content = sb.String()
		} else if it, ok := underlying.(*ir.InterfaceType); ok {
			var sb strings.Builder
			sb.WriteString("concept x\n")
			for _, m := range it.Methods {
				sb.WriteString(fmt.Sprintf("    x.%s() is %s\n", EscapeNimKeyword(m.Name), t.translateResults(m.Type.Results)))
			}
			content = sb.String()
		} else {
			content = types.MapType(underlying, t)
		}
		return &nim.TypeDecl{
			Name:       EscapeNimKeyword(d.Name),
			TypeParams: t.translateTypeParams(d.TypeParams),
			Content:    content,
			Exported:   isExported(d.Name),
			IsAlias:    d.Alias,
		}
	case *ir.ConstDecl:
		typ := types.MapType(d.Type, t)
		if strings.Contains(typ, "untyped") {
			typ = ""
		}
		return &nim.VarDecl{
			Kind:  "const",
			Name:  EscapeNimKeyword(d.Names[0]), // Simplified
			Typ:   typ,
			Value: t.translateExpr(d.Values[0]), // Simplified
		}
	case *ir.CGoDecl:
		return &nim.Stmt{Content: fmt.Sprintf("{.emit: \"\"\"\n%s\n\"\"\".}", d.Preamble)}
	case *ir.ImportDecl:
		if d.Path == "C" {
			return nil
		}
		// Map Go stdlib to gostdnim
		stdLibs := map[string]string{
			"builtin":       "gostdnim/builtin",
			"fmt":           "gostdnim/fmt",
			"os":            "gostdnim/os",
			"io":            "gostdnim/io",
			"errors":        "gostdnim/errors",
			"reflect":       "gostdnim/reflect",
			"sort":          "gostdnim/sort",
			"strconv":       "gostdnim/strconv",
			"strings":       "gostdnim/strings",
			"time":          "gostdnim/time",
			"unsafe":        "gostdnim/unsafe",
			"encoding/json": "gostdnim/encoding/json",
		}
		if nimPath, ok := stdLibs[d.Path]; ok {
			t.AddImport(nimPath)
			return nil
		}

		parts := strings.Split(d.Path, "/")
		pkg := parts[len(parts)-1]
		t.AddImport(pkg)
		return nil
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

func (t *translator) translateFields(fields []*ir.Field, variadic bool) []nim.Arg {
	var args []nim.Arg
	for i, f := range fields {
		for _, name := range f.Names {
			name = EscapeNimKeyword(name)
			typ := ""
			if variadic && i == len(fields)-1 {
				if st, ok := f.Type.(*ir.SliceType); ok {
					typ = "varargs[" + types.MapType(st.Elem, t) + "]"
				} else {
					typ = "varargs[" + types.MapType(f.Type, t) + "]"
				}
			} else {
				typ = types.MapType(f.Type, t)
			}
			args = append(args, nim.Arg{
				Name: name,
				Typ:  typ,
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
		res := types.MapType(results[0].Type, t)
		if res == "error" {
			res = "ref Exception"
		}
		return res
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
		if sn := t.translateStmt(stmt); sn != nil {
			nodes = append(nodes, sn)
		}
	}
	if len(nodes) == 0 {
		return []nim.Node{&nim.Stmt{Content: "discard"}}
	}
	return nodes
}

func (t *translator) translateStmt(stmt ir.Stmt) nim.Node {
	switch s := stmt.(type) {
	case *ir.ExprStmt:
		if ce, ok := s.X.(*ir.CallExpr); ok {
			if id, ok := ce.Fun.(*ir.Ident); ok && id.Name == "append" {
				args := ce.Args
				slice := t.translateExpr(args[0])
				val := t.translateExpr(args[1])
				return &nim.Stmt{Content: fmt.Sprintf("%s.add(%s)", slice, val)}
			}
		}
		return &nim.Stmt{Content: t.translateExprWithIndent(s.X, 0)}
	case *ir.ReturnStmt:
		if len(s.Results) == 0 {
			// Check for named returns
			return &nim.Stmt{Content: "return result"} // result is default in Nim
		}
		var res []string
		for _, r := range s.Results {
			res = append(res, t.translateExpr(r))
		}
		content := strings.Join(res, ", ")
		if len(res) > 1 {
			content = "(" + content + ")"
		}
		return &nim.Stmt{Content: "return " + content}
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
		if len(s.Lhs) == 1 && s.Lhs[0].(interface{String() string}).String() == "_" {
			return &nim.Stmt{Content: "discard " + t.translateExpr(s.Rhs[0])}
		}
		// Optimized append: s = append(s, x, y) -> s.add(x); s.add(y)
		if len(s.Lhs) == 1 && len(s.Rhs) == 1 && (s.Op == "=" || s.Op == ":=") {
			if ce, ok := s.Rhs[0].(*ir.CallExpr); ok {
				if id, ok := ce.Fun.(*ir.Ident); ok && id.Name == "append" {
					if len(ce.Args) >= 1 && t.translateExpr(s.Lhs[0]) == t.translateExpr(ce.Args[0]) {
						var adds []string
						if s.Op == ":=" {
							adds = append(adds, "var "+t.translateExpr(s.Lhs[0])+": "+types.MapType(s.Lhs[0].GetType(), t))
						}
						for i := 1; i < len(ce.Args); i++ {
							adds = append(adds, fmt.Sprintf("%s.add(%s)", t.translateExpr(s.Lhs[0]), t.translateExpr(ce.Args[i])))
						}
						if len(adds) == 0 {
							return &nim.Stmt{Content: "discard"}
						}
						return &nim.Stmt{Content: strings.Join(adds, "\n")}
					}
				}
			}
		}
		hasBlank := false
		for _, l := range s.Lhs {
			if l.(interface{String() string}).String() == "_" {
				hasBlank = true
				break
			}
		}
		if hasBlank && len(s.Lhs) > 1 {
			var rhs string
			if len(s.Rhs) == 1 {
				rhs = t.translateExpr(s.Rhs[0])
				// Ensure it's treated as a tuple for unpacking if needed
			} else {
				var rs []string
				for _, r := range s.Rhs {
					rs = append(rs, t.translateExpr(r))
				}
				rhs = "(" + strings.Join(rs, ", ") + ")"
			}
			var ls []string
			var discards []string
			for i, l := range s.Lhs {
				name := l.(interface{String() string}).String()
				if name == "_" {
					tmpName := fmt.Sprintf("_tmp%d", i)
					ls = append(ls, tmpName)
					discards = append(discards, "discard "+tmpName)
				} else {
					ls = append(ls, t.translateExpr(l))
				}
			}
			lhs := strings.Join(ls, ", ")
			op := s.Op
			if op == ":=" {
				return &nim.Stmt{Content: fmt.Sprintf("block:\n  var (%s) = %s\n  %s", lhs, rhs, strings.Join(discards, "\n  "))}
			}
			return &nim.Stmt{Content: fmt.Sprintf("block:\n  (%s) = %s\n  %s", lhs, rhs, strings.Join(discards, "\n  "))}
		}

		if len(s.Lhs) > 1 && len(s.Rhs) == 1 {
			if ta, ok := s.Rhs[0].(*ir.TypeAssertExpr); ok {
				// Comma-ok type assertion: v, ok := x.(T)
				lhs1 := t.translateExpr(s.Lhs[0])
				lhs2 := t.translateExpr(s.Lhs[1])
				typ := types.MapType(ta.Type, t)
				rhs := t.translateExpr(ta.X)
				if s.Op == ":=" {
					return &nim.Stmt{Content: fmt.Sprintf("var (%s, %s) = (if %s is %s: (%s(%s), true) else: (default(%s), false))", lhs1, lhs2, rhs, typ, typ, rhs, typ)}
				}
				return &nim.Stmt{Content: fmt.Sprintf("(%s, %s) = (if %s is %s: (%s(%s), true) else: (default(%s), false))", lhs1, lhs2, rhs, typ, typ, rhs, typ)}
			}
			// Case like: x, y := f() where f returns a tuple
			var lhs []string
			for _, l := range s.Lhs {
				lhs = append(lhs, t.translateExpr(l))
			}
			lhsStr := "(" + strings.Join(lhs, ", ") + ")"
			rhsStr := t.translateExpr(s.Rhs[0])
			if s.Op == ":=" {
				return &nim.Stmt{Content: "var " + lhsStr + " = " + rhsStr}
			}
			return &nim.Stmt{Content: lhsStr + " = " + rhsStr}
		}

		var lhs, rhs []string
		for _, l := range s.Lhs {
			lhs = append(lhs, t.translateExpr(l))
		}
		for _, r := range s.Rhs {
			rhs = append(rhs, t.translateExprWithIndent(r, 0))
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
			if len(lhs) > 1 {
				return &nim.Stmt{Content: "var " + lhsStr + " = " + rhsStr}
			}
			return &nim.Stmt{Content: "var " + lhsStr + ": " + types.MapType(s.Lhs[0].GetType(), t) + " = " + rhsStr}
		}
		return &nim.Stmt{Content: fmt.Sprintf("%s %s %s", lhsStr, op, rhsStr)}
	case *ir.BlockStmt:
		return &nim.BlockStmt{
			Body: t.translateBlock(s),
		}
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
		ifStmt := &nim.IfStmt{
			Cond: t.translateExpr(s.Cond),
			Body: t.translateBlock(s.Body),
			Else: elseNodes,
		}
		if s.Init != nil {
			return &nim.BlockStmt{
				Body: []nim.Node{
					t.translateStmt(s.Init),
					ifStmt,
				},
			}
		}
		return ifStmt
	case *ir.DeferStmt:
		return &nim.Stmt{Content: "defer:\n" + t.renderNodes([]nim.Node{&nim.Stmt{Content: t.translateExpr(s.Call)}}, 1)}
	case *ir.RangeStmt:
		key := "i"
		if s.Key != nil {
			key = t.translateExpr(s.Key)
			if key == "_" {
				key = "i"
			}
		}
		val := "v"
		if s.Value != nil {
			val = t.translateExpr(s.Value)
			if val == "_" {
				val = "v"
			}
		}
		return &nim.Stmt{Content: fmt.Sprintf("for %s, %s in %s:\n%s", key, val, t.translateExpr(s.X), t.renderNodes(t.translateBlock(s.Body), 1))}
	case *ir.SwitchStmt:
		cs := &nim.CaseStmt{Expr: t.translateExpr(s.Tag)}
		for _, stmt := range s.Body.List {
			tcc := stmt.(*ir.CaseClause)
			var vals []string
			for _, l := range tcc.List {
				vals = append(vals, t.translateExpr(l))
			}
			cs.Clauses = append(cs.Clauses, nim.CaseClause{
				Vals: vals,
				Body: t.translateBlock(&ir.BlockStmt{List: tcc.Body}),
			})
		}
		return cs
	case *ir.TypeSwitchStmt:
		// Map type switch to if-is chain
		var varName string
		var expr string
		if as, ok := s.Assign.(*ir.AssignStmt); ok {
			expr = t.translateExpr(as.Rhs[0])
			if len(as.Lhs) > 0 {
				varName = t.translateExpr(as.Lhs[0])
			}
		} else if es, ok := s.Assign.(*ir.ExprStmt); ok {
			// Type switch like switch any.(type)
			if ta, ok := es.X.(*ir.TypeAssertExpr); ok {
				expr = t.translateExpr(ta.X)
			} else {
				expr = t.translateExpr(es.X)
			}
		}

		var first *nim.IfStmt
		var current *nim.IfStmt
		for _, stmt := range s.Body.List {
			tcc := stmt.(*ir.CaseClause)
			if len(tcc.List) > 0 {
				typStr := t.translateExpr(tcc.List[0])
				cond := fmt.Sprintf("%s is %s", expr, typStr)
				body := t.translateBlock(&ir.BlockStmt{List: tcc.Body})
				if varName != "" && varName != "_" {
					body = append([]nim.Node{&nim.Stmt{Content: fmt.Sprintf("let %s = %s(%s)", varName, typStr, expr)}}, body...)
				}
				branch := &nim.IfStmt{
					Cond: cond,
					Body: body,
				}
				if first == nil {
					first = branch
				} else {
					current.Else = []nim.Node{branch}
				}
				current = branch
			} else {
				// default case
				if current != nil {
					current.Else = t.translateBlock(&ir.BlockStmt{List: tcc.Body})
				}
			}
		}
		if first == nil {
			return &nim.Stmt{Content: "discard # empty type switch"}
		}
		return first
	case *ir.CaseClause:
		var labels []string
		for _, l := range s.List {
			labels = append(labels, t.translateExpr(l))
		}
		label := strings.Join(labels, ", ")
		if label == "" {
			label = "else"
		} else {
			label = "of " + label
		}
		return &nim.Stmt{Content: fmt.Sprintf("%s:\n%s", label, t.renderNodes(t.translateBlock(&ir.BlockStmt{List: s.Body}), 1))}
	case *ir.ForStmt:
		cond := t.translateExpr(s.Cond)
		body := t.translateBlock(s.Body)
		if s.Post != nil {
			body = append(body, t.translateStmt(s.Post))
		}
		forStmt := &nim.ForStmt{
			Cond: cond,
			Body: body,
		}
		if s.Init != nil {
			return &nim.BlockStmt{
				Body: []nim.Node{
					t.translateStmt(s.Init),
					forStmt,
				},
			}
		}
		return forStmt
	case *ir.IncDecStmt:
		return &nim.Stmt{Content: fmt.Sprintf("%s.inc", t.translateExpr(s.X))}
	case *ir.BranchStmt:
		if s.Label != "" {
			return &nim.Stmt{Content: strings.ToLower(s.Tok) + " " + s.Label}
		}
		return &nim.Stmt{Content: strings.ToLower(s.Tok)}
	}
	return &nim.Stmt{Content: "discard # unknown stmt"}
}

func (t *translator) renderNodes(nodes []nim.Node, indent int) string {
	var sb strings.Builder
	if len(nodes) == 0 {
		return strings.Repeat("  ", indent) + "discard"
	}
	for i, n := range nodes {
		if n == nil {
			continue
		}
		if i > 0 {
			sb.WriteString("\n")
		}
		sb.WriteString(n.Render(indent))
	}
	if sb.Len() == 0 {
		return strings.Repeat("  ", indent) + "discard"
	}
	return sb.String()
}

func (t *translator) isTypeName(expr ir.Expr) bool {
	if id, ok := expr.(*ir.Ident); ok {
		return isExported(id.Name) // Basic heuristic
	}
	return false
}

func (t *translator) translateExpr(expr ir.Expr) string {
	return t.translateExprWithIndent(expr, 0)
}

func EscapeNimKeyword(name string) string {
	keywords := map[string]bool{
		"addr":      true,
		"and":       true,
		"as":        true,
		"asm":       true,
		"bind":      true,
		"block":     true,
		"break":     true,
		"case":      true,
		"cast":      true,
		"concept":   true,
		"const":     true,
		"continue":  true,
		"converter": true,
		"defer":     true,
		"discard":   true,
		"distinct":  true,
		"div":       true,
		"do":        true,
		"elif":      true,
		"else":      true,
		"end":       true,
		"enum":      true,
		"except":    true,
		"export":    true,
		"finally":   true,
		"for":       true,
		"from":      true,
		"func":      true,
		"if":        true,
		"import":    true,
		"in":        true,
		"include":   true,
		"interface": true,
		"is":        true,
		"isnot":     true,
		"iterator":  true,
		"let":       true,
		"macro":     true,
		"method":    true,
		"mixin":     true,
		"mod":       true,
		"not":       true,
		"notin":     true,
		"object":    true,
		"of":        true,
		"or":        true,
		"out":       true,
		"proc":      true,
		"ptr":       true,
		"ptr_":      true,
		"raise":     true,
		"ref":       true,
		"return":    true,
		"shl":       true,
		"shr":       true,
		"static":    true,
		"template":  true,
		"try":       true,
		"tuple":     true,
		"type":      true,
		"using":     true,
		"var":       true,
		"when":      true,
		"while":     true,
		"xor":       true,
		"yield":     true,
	}
	if keywords[name] {
		return name + "_"
	}
	return name
}

func (t *translator) translateExprWithIndent(expr ir.Expr, n int) string {
	if expr == nil {
		return ""
	}
	switch e := expr.(type) {
	case *ir.Ident:
		return EscapeNimKeyword(e.Name)
	case *ir.BasicLit:
		return e.Value
	case *ir.CallExpr:
		fun := t.translateExpr(e.Fun)
		if fun == "make" && len(e.Args) >= 1 {
			if te, ok := e.Args[0].(*ir.TypeExpr); ok {
				switch mt := te.Type.(type) {
				case *ir.MapType:
					t.AddImport("tables")
					return "initTable[" + types.MapType(mt.Key, t) + ", " + types.MapType(mt.Value, t) + "]()"
				case *ir.SliceType:
					size := "0"
					if len(e.Args) >= 2 {
						size = t.translateExpr(e.Args[1])
					}
					return "newSeq[" + types.MapType(mt.Elem, t) + "](" + size + ")"
				}
			}
		}

		isMember := false
		if b, ok := builtins.Builtins[fun]; ok {
			fun = b.NimName
			isMember = b.IsMember
		}

		if fun == "panic" {
			arg := t.translateExpr(e.Args[0])
			return fmt.Sprintf("raise newException(Exception, %s)", arg)
		}
		if fun == "recover" {
			return "getCurrentExceptionMsg()" // Simplified, Go's recover is more complex
		}

		var args []string
		for _, arg := range e.Args {
			args = append(args, t.translateExprWithIndent(arg, n+1))
		}
		if isMember && len(args) > 0 {
			if fun == "add" {
				// Go's append returns the slice, Nim's add is void.
				// For now, if it's used as an expression, we need to handle it.
				// This is a common transpilation challenge.
			if len(args) > 2 {
				var sb strings.Builder
				sb.WriteString("(var temp = ")
				sb.WriteString(args[0])
				sb.WriteString(";")
				for i := 1; i < len(args); i++ {
					sb.WriteString(" temp.add(")
					sb.WriteString(args[i])
					sb.WriteString(");")
				}
				sb.WriteString(" temp)")
				return sb.String()
			}
				return fmt.Sprintf("(var temp = %s; temp.add(%s); temp)", args[0], strings.Join(args[1:], ", "))
			}
			return fmt.Sprintf("%s.%s(%s)", args[0], fun, strings.Join(args[1:], ", "))
		}
		return fmt.Sprintf("%s(%s)", fun, strings.Join(args, ", "))
	case *ir.TypeAssertExpr:
		typ := types.MapType(e.Type, t)
		if typ == "any" {
			return t.translateExpr(e.X)
		}
		// Single-value type assertion: x.(T) - should panic in Go if it fails
		return fmt.Sprintf("(if %s is %s: %s(%s) else: (raise newException(Exception, \"type assertion failed\"); default(%s)))", t.translateExpr(e.X), typ, typ, t.translateExpr(e.X), typ)
	case *ir.FuncLit:
		var params []string
		for _, p := range e.Type.Params {
			for _, name := range p.Names {
				params = append(params, fmt.Sprintf("%s: %s", name, types.MapType(p.Type, t)))
			}
		}
		ret := "void"
		if len(e.Type.Results) > 0 {
			ret = types.MapType(e.Type.Results[0].Type, t)
		}
		// Use a single-line proc if the body is simple, otherwise multiline
		bodyNodes := t.translateBlock(e.Body)
		if len(bodyNodes) == 1 {
			if stmt, ok := bodyNodes[0].(*nim.Stmt); ok && !strings.Contains(stmt.Content, "\n") {
				content := stmt.Content
				if strings.HasPrefix(content, "return ") {
					content = strings.TrimPrefix(content, "return ")
				}
				return fmt.Sprintf("(proc(%s): %s = %s)", strings.Join(params, ", "), ret, content)
			}
		}
		var sb strings.Builder
		for i, bn := range bodyNodes {
			if i > 0 {
				sb.WriteString("\n")
			}
			res := bn.Render(n + 1)
			if !strings.HasPrefix(res, "  ") {
				res = strings.Repeat("  ", n+1) + res
			}
			sb.WriteString(res)
		}
		return fmt.Sprintf("(proc(%s): %s =\n%s)", strings.Join(params, ", "), ret, sb.String())
	case *ir.BinaryExpr:
		if e.Op == ":" {
			return fmt.Sprintf("%s: %s", t.translateExpr(e.X), t.translateExpr(e.Y))
		}
		op := e.Op
		// Map Go operators to Nim if they differ
		switch op {
		case "==":
			if t.translateExpr(e.Y) == "nil" {
				return fmt.Sprintf("%s.isNil", t.translateExpr(e.X))
			}
		case "!=":
			if t.translateExpr(e.Y) == "nil" {
				return fmt.Sprintf("not %s.isNil", t.translateExpr(e.X))
			}
		case "+":
			if e.X.GetType().String() == "string" || e.Y.GetType().String() == "string" {
				op = "&"
			}
		case "&&":
			op = "and"
		case "||":
			op = "or"
		case "!":
			op = "not" // though this is Unary
		case ":":
			return fmt.Sprintf("%s: %s", t.translateExpr(e.X), t.translateExpr(e.Y))
		}
		return fmt.Sprintf("(%s %s %s)", t.translateExpr(e.X), op, t.translateExpr(e.Y))
	case *ir.SelectorExpr:
		if x, ok := e.X.(*ir.Ident); ok && x.Name == "C" {
			return e.Sel
		}
		if t.isTypeName(e.X) {
			// Method expression Type.Method
			return t.translateExpr(e.X) + "." + e.Sel
		}
		return fmt.Sprintf("%s.%s", t.translateExpr(e.X), e.Sel)
	case *ir.IndexExpr:
		return fmt.Sprintf("%s[%s]", t.translateExpr(e.X), t.translateExpr(e.Index))
	case *ir.UnaryExpr:
		op := e.Op
		switch op {
		case "!":
			op = "not "
		case "*":
			return fmt.Sprintf("%s[]", t.translateExpr(e.X))
		case "&":
			return fmt.Sprintf("addr(%s)", t.translateExpr(e.X))
		}
		return fmt.Sprintf("%s%s", op, t.translateExpr(e.X))
	case *ir.CompositeLit:
		var elms []string
		for _, elm := range e.Elms {
			elms = append(elms, t.translateExpr(elm))
		}
		typ := types.MapType(e.Type, t)
		if strings.HasPrefix(typ, "seq") {
			return fmt.Sprintf("@ [%s]", strings.Join(elms, ", "))
		}
		if strings.HasPrefix(typ, "Table") {
			return fmt.Sprintf("{%s}.toTable", strings.Join(elms, ", "))
		}
		if strings.HasPrefix(typ, "array") {
			return fmt.Sprintf("[%s]", strings.Join(elms, ", "))
		}
		if typ == "object" || strings.HasPrefix(typ, "struct") || strings.HasPrefix(typ, "tuple") || isExported(typ) || typ == "Person" || typ == "Employee" || typ == "MyInt" || typ == "MyError" || typ == "Point" || strings.Contains(typ, ".") || typ == "AliasInt" {
			return fmt.Sprintf("%s(%s)", typ, strings.Join(elms, ", "))
		}
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
