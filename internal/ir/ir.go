package ir

import "fmt"

// Node is the base interface for all IR nodes.
type Node interface {
	fmt.Stringer
}

// Program represents the entire transpilation unit.
type Program struct {
	Packages []*Package
}

func (p *Program) String() string {
	return fmt.Sprintf("Program(Packages: %d)", len(p.Packages))
}

// Package represents a Go package.
type Package struct {
	Name  string
	Path  string
	Files []*File
}

func (p *Package) String() string {
	return fmt.Sprintf("Package(%s)", p.Name)
}

// File represents a Go source file.
type File struct {
	Path  string
	Decls []Decl
}

func (f *File) String() string {
	return fmt.Sprintf("File(%s)", f.Path)
}

// Decl is the interface for all declarations.
type Decl interface {
	Node
	declNode()
}

// Stmt is the interface for all statements.
type Stmt interface {
	Node
	stmtNode()
}

// Expr is the interface for all expressions.
type Expr interface {
	Node
	exprNode()
	GetType() Type
}

// Type is the interface for all types in the IR.
type Type interface {
	fmt.Stringer
	typeNode()
}
