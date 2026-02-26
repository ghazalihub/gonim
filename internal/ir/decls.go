package ir

import "fmt"

type (
	FuncDecl struct {
		Name     string
		Receiver *Field
		Type     *FuncType
		Body     *BlockStmt
	}

	TypeDecl struct {
		Name       string
		Type       Type
		Alias      bool
		TypeParams []*TypeParam
	}

	VarDecl struct {
		Names  []string
		Type   Type
		Values []Expr
		Embeds []string
	}

	ConstDecl struct {
		Names  []string
		Type   Type
		Values []Expr
	}

	ImportDecl struct {
		Path string
		Name string
	}

	CGoDecl struct {
		Preamble string
	}
)

func (*FuncDecl) declNode()   {}
func (*TypeDecl) declNode()   {}
func (*VarDecl) declNode()    {}
func (*ConstDecl) declNode()  {}
func (*ImportDecl) declNode() {}
func (*CGoDecl) declNode()    {}

func (d *FuncDecl) String() string   { return fmt.Sprintf("func %s", d.Name) }
func (d *TypeDecl) String() string   { return fmt.Sprintf("type %s", d.Name) }
func (d *VarDecl) String() string    { return fmt.Sprintf("var %v", d.Names) }
func (d *ConstDecl) String() string  { return fmt.Sprintf("const %v", d.Names) }
func (d *ImportDecl) String() string { return fmt.Sprintf("import %s", d.Path) }
func (d *CGoDecl) String() string    { return "cgo preamble" }
