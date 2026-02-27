package ir

import (
	"fmt"
)

type (
	BasicType struct {
		Name string
	}

	StructType struct {
		Fields []*Field
	}

	Field struct {
		Names []string
		Type  Type
		Tag   string
	}

	InterfaceType struct {
		Methods []*FuncDecl
	}

	SliceType struct {
		Elem Type
	}

	ArrayType struct {
		Len  int64
		Elem Type
	}

	MapType struct {
		Key  Type
		Value Type
	}

	PointerType struct {
		Elem Type
	}

	ChanType struct {
		Dir  int // 0: both, 1: send, 2: recv
		Elem Type
	}

	FuncType struct {
		Params  []*Field
		Results []*Field
		TypeParams []*TypeParam
		Variadic   bool
	}

	TypeParam struct {
		Name string
		Constraint Type
	}

	NamedType struct {
		Package string
		Name    string
		Underlying Type
		TypeArgs   []Type
	}
)

func (*BasicType) typeNode()     {}
func (*StructType) typeNode()    {}
func (*InterfaceType) typeNode() {}
func (*SliceType) typeNode()     {}
func (*ArrayType) typeNode()     {}
func (*MapType) typeNode()       {}
func (*PointerType) typeNode()   {}
func (*ChanType) typeNode()      {}
func (*FuncType) typeNode()      {}
func (*TypeParam) typeNode()     {}
func (*NamedType) typeNode()     {}

func (t *BasicType) String() string { return t.Name }
func (t *StructType) String() string { return "struct{...}" }
func (t *InterfaceType) String() string { return "interface{...}" }
func (t *SliceType) String() string { return fmt.Sprintf("[]%s", t.Elem) }
func (t *ArrayType) String() string { return fmt.Sprintf("[%d]%s", t.Len, t.Elem) }
func (t *MapType) String() string { return fmt.Sprintf("map[%s]%s", t.Key, t.Value) }
func (t *PointerType) String() string { return fmt.Sprintf("*%s", t.Elem) }
func (t *ChanType) String() string    { return fmt.Sprintf("chan %s", t.Elem) }
func (t *FuncType) String() string    { return "func(...)..." }
func (t *TypeParam) String() string { return t.Name }
func (t *NamedType) String() string {
	if t.Package != "" {
		return fmt.Sprintf("%s.%s", t.Package, t.Name)
	}
	return t.Name
}
