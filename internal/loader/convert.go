package loader

import (
	"go/ast"
	"go/token"
	"go/types"
	"strings"

	"github.com/user/go2nim/internal/ir"
	"golang.org/x/tools/go/packages"
)

type converter struct {
	pkg     *packages.Package
	typeMap map[types.Type]ir.Type
}

func convertPackage(pkg *packages.Package) (*ir.Package, error) {
	c := &converter{
		pkg:     pkg,
		typeMap: make(map[types.Type]ir.Type),
	}

	irPkg := &ir.Package{
		Name: pkg.Name,
		Path: pkg.PkgPath,
	}

	for _, file := range pkg.Syntax {
		irFile, err := c.convertFile(file)
		if err != nil {
			return nil, err
		}
		irPkg.Files = append(irPkg.Files, irFile)
	}

	return irPkg, nil
}

func (c *converter) convertFile(f *ast.File) (*ir.File, error) {
	irFile := &ir.File{
		Path: c.pkg.Fset.File(f.Pos()).Name(),
	}

	for _, decl := range f.Decls {
		irDecls := c.convertDecl(decl)
		irFile.Decls = append(irFile.Decls, irDecls...)
	}

	return irFile, nil
}

func (c *converter) convertDecl(decl ast.Decl) []ir.Decl {
	var decls []ir.Decl
	switch d := decl.(type) {
	case *ast.GenDecl:
		switch d.Tok {
		case token.IMPORT:
			for _, s := range d.Specs {
				spec := s.(*ast.ImportSpec)
				if spec.Path.Value == "\"C\"" {
					preamble := ""
					if d.Doc != nil {
						preamble = d.Doc.Text()
					}
					decls = append(decls, &ir.CGoDecl{Preamble: preamble})
					continue
				}
				name := ""
				if spec.Name != nil {
					name = spec.Name.Name
				}
				decls = append(decls, &ir.ImportDecl{
					Path: strings.Trim(spec.Path.Value, "\""),
					Name: name,
				})
			}
		case token.TYPE:
			for _, s := range d.Specs {
				spec := s.(*ast.TypeSpec)
				typ := c.pkg.TypesInfo.Defs[spec.Name].Type()
				var tps []*ir.TypeParam
				if named, ok := typ.(*types.Named); ok && named.TypeParams() != nil {
					for i := 0; i < named.TypeParams().Len(); i++ {
						tp := named.TypeParams().At(i)
						tps = append(tps, &ir.TypeParam{
							Name:       tp.Obj().Name(),
							Constraint: c.convertType(tp.Constraint()),
						})
					}
				}
				decls = append(decls, &ir.TypeDecl{
					Name:       spec.Name.Name,
					Type:       c.convertType(typ),
					TypeParams: tps,
				})
			}
		case token.VAR:
			for _, s := range d.Specs {
				spec := s.(*ast.ValueSpec)
				var names []string
				for _, name := range spec.Names {
					names = append(names, name.Name)
				}
				var values []ir.Expr
				for _, val := range spec.Values {
					values = append(values, c.convertExpr(val))
				}
				var embeds []string
				if d.Doc != nil {
					for _, comment := range d.Doc.List {
						if strings.HasPrefix(comment.Text, "//go:embed") {
							embeds = append(embeds, strings.Fields(strings.TrimPrefix(comment.Text, "//go:embed"))...)
						}
					}
				}
				decls = append(decls, &ir.VarDecl{
					Names:  names,
					Type:   c.convertType(c.pkg.TypesInfo.Defs[spec.Names[0]].Type()),
					Values: values,
					Embeds: embeds,
				})
			}
		case token.CONST:
			for _, s := range d.Specs {
				spec := s.(*ast.ValueSpec)
				var names []string
				var values []ir.Expr
				for _, name := range spec.Names {
					names = append(names, name.Name)
					if obj, ok := c.pkg.TypesInfo.Defs[name].(*types.Const); ok {
						val := obj.Val()
						values = append(values, &ir.BasicLit{
							Value: val.ExactString(),
							Kind:  val.Kind().String(),
							Typ:   c.convertType(obj.Type()),
						})
					}
				}
				decls = append(decls, &ir.ConstDecl{
					Names:  names,
					Type:   c.convertType(c.pkg.TypesInfo.Defs[spec.Names[0]].Type()),
					Values: values,
				})
			}
		}
	case *ast.FuncDecl:
		irFunc := &ir.FuncDecl{
			Name: d.Name.Name,
			Type: c.convertType(c.pkg.TypesInfo.Defs[d.Name].Type()).(*ir.FuncType),
			Body: c.convertBlockStmt(d.Body),
		}
		if d.Recv != nil {
			irFunc.Receiver = c.convertField(d.Recv.List[0])[0]
		}
		decls = append(decls, irFunc)
	}
	return decls
}

func (c *converter) convertType(t types.Type) ir.Type {
	if it, ok := c.typeMap[t]; ok {
		return it
	}

	var irType ir.Type
	switch tt := t.(type) {
	case *types.Basic:
		irType = &ir.BasicType{Name: tt.Name()}
	case *types.Pointer:
		irType = &ir.PointerType{Elem: c.convertType(tt.Elem())}
	case *types.Slice:
		irType = &ir.SliceType{Elem: c.convertType(tt.Elem())}
	case *types.Array:
		irType = &ir.ArrayType{Len: tt.Len(), Elem: c.convertType(tt.Elem())}
	case *types.Map:
		irType = &ir.MapType{Key: c.convertType(tt.Key()), Value: c.convertType(tt.Elem())}
	case *types.Struct:
		fields := make([]*ir.Field, tt.NumFields())
		for i := 0; i < tt.NumFields(); i++ {
			f := tt.Field(i)
			fields[i] = &ir.Field{
				Names: []string{f.Name()},
				Type:  c.convertType(f.Type()),
				Tag:   tt.Tag(i),
			}
		}
		irType = &ir.StructType{Fields: fields}
	case *types.Interface:
		methods := make([]*ir.FuncDecl, tt.NumMethods())
		for i := 0; i < tt.NumMethods(); i++ {
			m := tt.Method(i)
			methods[i] = &ir.FuncDecl{
				Name: m.Name(),
				Type: c.convertType(m.Type()).(*ir.FuncType),
			}
		}
		irType = &ir.InterfaceType{Methods: methods}
	case *types.Signature:
		irType = c.convertSignature(tt)
	case *types.Named:
		obj := tt.Obj()
		pkgName := ""
		if obj.Pkg() != nil {
			pkgName = obj.Pkg().Name()
		}
		nt := &ir.NamedType{
			Package: pkgName,
			Name:    obj.Name(),
		}
		if tt.TypeArgs() != nil {
			for i := 0; i < tt.TypeArgs().Len(); i++ {
				nt.TypeArgs = append(nt.TypeArgs, c.convertType(tt.TypeArgs().At(i)))
			}
		}
		c.typeMap[t] = nt // Register before converting underlying to avoid recursion
		nt.Underlying = c.convertType(tt.Underlying())
		irType = nt
	case *types.TypeParam:
		irType = &ir.TypeParam{Name: tt.Obj().Name()}
	default:
		irType = &ir.BasicType{Name: "unknown"}
	}

	c.typeMap[t] = irType
	return irType
}

func (c *converter) convertSignature(sig *types.Signature) *ir.FuncType {
	ft := &ir.FuncType{Variadic: sig.Variadic()}
	if sig.TypeParams() != nil {
		for i := 0; i < sig.TypeParams().Len(); i++ {
			tp := sig.TypeParams().At(i)
			ft.TypeParams = append(ft.TypeParams, &ir.TypeParam{
				Name:       tp.Obj().Name(),
				Constraint: c.convertType(tp.Constraint()),
			})
		}
	}
	for i := 0; i < sig.Params().Len(); i++ {
		p := sig.Params().At(i)
		ft.Params = append(ft.Params, &ir.Field{
			Names: []string{p.Name()},
			Type:  c.convertType(p.Type()),
		})
	}
	for i := 0; i < sig.Results().Len(); i++ {
		r := sig.Results().At(i)
		ft.Results = append(ft.Results, &ir.Field{
			Names: []string{r.Name()},
			Type:  c.convertType(r.Type()),
		})
	}
	// TODO: TypeParams
	return ft
}

func (c *converter) convertField(f *ast.Field) []*ir.Field {
	irType := c.convertType(c.pkg.TypesInfo.TypeOf(f.Type))
	var fields []*ir.Field
	if len(f.Names) == 0 {
		fields = append(fields, &ir.Field{Type: irType})
	} else {
		for _, name := range f.Names {
			fields = append(fields, &ir.Field{
				Names: []string{name.Name},
				Type:  irType,
			})
		}
	}
	return fields
}

func (c *converter) convertBlockStmt(b *ast.BlockStmt) *ir.BlockStmt {
	if b == nil {
		return nil
	}
	irBlock := &ir.BlockStmt{}
	for _, stmt := range b.List {
		irStmt := c.convertStmt(stmt)
		if irStmt != nil {
			irBlock.List = append(irBlock.List, irStmt)
		}
	}
	return irBlock
}

func (c *converter) convertStmt(stmt ast.Stmt) ir.Stmt {
	switch s := stmt.(type) {
	case *ast.ExprStmt:
		return &ir.ExprStmt{X: c.convertExpr(s.X)}
	case *ast.DeclStmt:
		return &ir.DeclStmt{Decls: c.convertDecl(s.Decl)}
	case *ast.AssignStmt:
		as := &ir.AssignStmt{Op: s.Tok.String()}
		for _, lhs := range s.Lhs {
			as.Lhs = append(as.Lhs, c.convertExpr(lhs))
		}
		for _, rhs := range s.Rhs {
			as.Rhs = append(as.Rhs, c.convertExpr(rhs))
		}
		return as
	case *ast.ReturnStmt:
		rs := &ir.ReturnStmt{}
		for _, res := range s.Results {
			rs.Results = append(rs.Results, c.convertExpr(res))
		}
		return rs
	case *ast.IfStmt:
		return &ir.IfStmt{
			Init: c.convertStmt(s.Init),
			Cond: c.convertExpr(s.Cond),
			Body: c.convertBlockStmt(s.Body),
			Else: c.convertStmt(s.Else),
		}
	case *ast.ForStmt:
		return &ir.ForStmt{
			Init: c.convertStmt(s.Init),
			Cond: c.convertExpr(s.Cond),
			Post: c.convertStmt(s.Post),
			Body: c.convertBlockStmt(s.Body),
		}
	case *ast.RangeStmt:
		return &ir.RangeStmt{
			Key:   c.convertExpr(s.Key),
			Value: c.convertExpr(s.Value),
			X:     c.convertExpr(s.X),
			Body:  c.convertBlockStmt(s.Body),
		}
	case *ast.BranchStmt:
		return &ir.BranchStmt{
			Tok:   s.Tok.String(),
			Label: s.Label.Name,
		}
	case *ast.DeferStmt:
		return &ir.DeferStmt{
			Call: c.convertExpr(s.Call).(*ir.CallExpr),
		}
	case *ast.GoStmt:
		return &ir.GoStmt{
			Call: c.convertExpr(s.Call).(*ir.CallExpr),
		}
	case *ast.SwitchStmt:
		return &ir.SwitchStmt{
			Init: c.convertStmt(s.Init),
			Tag:  c.convertExpr(s.Tag),
			Body: c.convertBlockStmt(s.Body),
		}
	case *ast.CaseClause:
		cc := &ir.CaseClause{}
		for _, expr := range s.List {
			cc.List = append(cc.List, c.convertExpr(expr))
		}
		for _, stmt := range s.Body {
			cc.Body = append(cc.Body, c.convertStmt(stmt))
		}
		return cc
	}
	return nil
}

func (c *converter) convertExpr(expr ast.Expr) ir.Expr {
	switch e := expr.(type) {
	case *ast.Ident:
		return &ir.Ident{
			Name: e.Name,
			Typ:  c.convertType(c.pkg.TypesInfo.TypeOf(e)),
		}
	case *ast.BasicLit:
		return &ir.BasicLit{
			Value: e.Value,
			Kind:  e.Kind.String(),
			Typ:   c.convertType(c.pkg.TypesInfo.TypeOf(e)),
		}
	case *ast.BinaryExpr:
		return &ir.BinaryExpr{
			X:   c.convertExpr(e.X),
			Op:  e.Op.String(),
			Y:   c.convertExpr(e.Y),
			Typ: c.convertType(c.pkg.TypesInfo.TypeOf(e)),
		}
	case *ast.CallExpr:
		ce := &ir.CallExpr{
			Fun: c.convertExpr(e.Fun),
			Typ: c.convertType(c.pkg.TypesInfo.TypeOf(e)),
		}
		for _, arg := range e.Args {
			ce.Args = append(ce.Args, c.convertExpr(arg))
		}
		return ce
	case *ast.SelectorExpr:
		return &ir.SelectorExpr{
			X:   c.convertExpr(e.X),
			Sel: e.Sel.Name,
			Typ: c.convertType(c.pkg.TypesInfo.TypeOf(e)),
		}
	case *ast.UnaryExpr:
		return &ir.UnaryExpr{
			Op:  e.Op.String(),
			X:   c.convertExpr(e.X),
			Typ: c.convertType(c.pkg.TypesInfo.TypeOf(e)),
		}
	case *ast.CompositeLit:
		cl := &ir.CompositeLit{
			Type: c.convertType(c.pkg.TypesInfo.TypeOf(e)),
		}
		for _, elm := range e.Elts {
			cl.Elms = append(cl.Elms, c.convertExpr(elm))
		}
		return cl
	case *ast.IndexExpr:
		return &ir.IndexExpr{
			X:     c.convertExpr(e.X),
			Index: c.convertExpr(e.Index),
			Typ:   c.convertType(c.pkg.TypesInfo.TypeOf(e)),
		}
	case *ast.IndexListExpr:
		// Go 1.18+ multiple type args
		var indices []ir.Expr
		for _, idx := range e.Indices {
			indices = append(indices, c.convertExpr(idx))
		}
		// Using IndexExpr for now or maybe a new IR type?
		// Let's just use the first one for simplicity or create a new node.
		return &ir.IndexExpr{
			X:     c.convertExpr(e.X),
			Index: indices[0], // Simplified
			Typ:   c.convertType(c.pkg.TypesInfo.TypeOf(e)),
		}
	case *ast.SliceExpr:
		return &ir.SliceExpr{
			X:    c.convertExpr(e.X),
			Low:  c.convertExpr(e.Low),
			High: c.convertExpr(e.High),
			Max:  c.convertExpr(e.Max),
			Typ:  c.convertType(c.pkg.TypesInfo.TypeOf(e)),
		}
	case *ast.TypeAssertExpr:
		return &ir.TypeAssertExpr{
			X:    c.convertExpr(e.X),
			Type: c.convertType(c.pkg.TypesInfo.TypeOf(e.Type)),
			Typ:  c.convertType(c.pkg.TypesInfo.TypeOf(e)),
		}
	case *ast.FuncLit:
		return &ir.FuncLit{
			Type: c.convertType(c.pkg.TypesInfo.TypeOf(e.Type)).(*ir.FuncType),
			Body: c.convertBlockStmt(e.Body),
		}
	case *ast.ArrayType, *ast.MapType, *ast.ChanType, *ast.InterfaceType, *ast.StructType:
		return &ir.TypeExpr{
			Type: c.convertType(c.pkg.TypesInfo.TypeOf(e.(ast.Expr))),
		}
	case *ast.KeyValueExpr:
		return &ir.BinaryExpr{
			X:   c.convertExpr(e.Key),
			Op:  ":",
			Y:   c.convertExpr(e.Value),
			Typ: c.convertType(c.pkg.TypesInfo.TypeOf(e)),
		}
	}
	return &ir.Ident{Name: "unknown_expr"}
}
