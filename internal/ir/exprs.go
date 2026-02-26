package ir

import "fmt"

type (
	Ident struct {
		Name string
		Typ  Type
	}

	BasicLit struct {
		Value string
		Kind  string // INT, FLOAT, IMAG, CHAR, STRING
		Typ   Type
	}

	CompositeLit struct {
		Type Type
		Elms []Expr
	}

	UnaryExpr struct {
		Op string
		X  Expr
		Typ Type
	}

	BinaryExpr struct {
		X  Expr
		Op string
		Y  Expr
		Typ Type
	}

	CallExpr struct {
		Fun  Expr
		Args []Expr
		Typ  Type
	}

	SelectorExpr struct {
		X   Expr
		Sel string
		Typ Type
	}

	IndexExpr struct {
		X     Expr
		Index Expr
		Typ   Type
	}

	SliceExpr struct {
		X    Expr
		Low  Expr
		High Expr
		Max  Expr
		Typ  Type
	}

	TypeAssertExpr struct {
		X    Expr
		Type Type
		Typ  Type
	}

	FuncLit struct {
		Type *FuncType
		Body *BlockStmt
	}

	TypeExpr struct {
		Type Type
	}
)

func (*Ident) exprNode()          {}
func (*BasicLit) exprNode()       {}
func (*CompositeLit) exprNode()   {}
func (*UnaryExpr) exprNode()      {}
func (*BinaryExpr) exprNode()     {}
func (*CallExpr) exprNode()       {}
func (*SelectorExpr) exprNode()   {}
func (*IndexExpr) exprNode()      {}
func (*SliceExpr) exprNode()      {}
func (*TypeAssertExpr) exprNode() {}
func (*FuncLit) exprNode()        {}
func (*TypeExpr) exprNode()       {}

func (e *Ident) GetType() Type          { return e.Typ }
func (e *BasicLit) GetType() Type       { return e.Typ }
func (e *CompositeLit) GetType() Type   { return e.Type }
func (e *UnaryExpr) GetType() Type      { return e.Typ }
func (e *BinaryExpr) GetType() Type     { return e.Typ }
func (e *CallExpr) GetType() Type       { return e.Typ }
func (e *SelectorExpr) GetType() Type   { return e.Typ }
func (e *IndexExpr) GetType() Type      { return e.Typ }
func (e *SliceExpr) GetType() Type      { return e.Typ }
func (e *TypeAssertExpr) GetType() Type { return e.Typ }
func (e *FuncLit) GetType() Type        { return e.Type }
func (e *TypeExpr) GetType() Type       { return e.Type }

func (e *Ident) String() string          { return e.Name }
func (e *BasicLit) String() string       { return e.Value }
func (e *CompositeLit) String() string   { return "composite lit" }
func (e *UnaryExpr) String() string      { return fmt.Sprintf("(%s%s)", e.Op, e.X) }
func (e *BinaryExpr) String() string     { return fmt.Sprintf("(%s %s %s)", e.X, e.Op, e.Y) }
func (e *CallExpr) String() string       { return fmt.Sprintf("%s(...)", e.Fun) }
func (e *SelectorExpr) String() string   { return fmt.Sprintf("%s.%s", e.X, e.Sel) }
func (e *IndexExpr) String() string      { return fmt.Sprintf("%s[%s]", e.X, e.Index) }
func (e *SliceExpr) String() string      { return "slice expr" }
func (e *TypeAssertExpr) String() string { return "type assert" }
func (e *FuncLit) String() string        { return "func lit" }
func (e *TypeExpr) String() string       { return fmt.Sprintf("type(%s)", e.Type) }
