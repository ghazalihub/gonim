package ir

import "fmt"

type (
	BlockStmt struct {
		List []Stmt
	}

	ExprStmt struct {
		X Expr
	}

	DeclStmt struct {
		Decls []Decl
	}

	AssignStmt struct {
		Lhs []Expr
		Rhs []Expr
		Op  string // e.g., "=", ":=", "+="
	}

	ReturnStmt struct {
		Results []Expr
	}

	IfStmt struct {
		Init Stmt
		Cond Expr
		Body *BlockStmt
		Else Stmt
	}

	ForStmt struct {
		Init Stmt
		Cond Expr
		Post Stmt
		Body *BlockStmt
	}

	RangeStmt struct {
		Key   Expr
		Value Expr
		X     Expr
		Body  *BlockStmt
	}

	SendStmt struct {
		Chan  Expr
		Value Expr
	}

	LabeledStmt struct {
		Label string
		Stmt  Stmt
	}

	BranchStmt struct {
		Tok   string // BREAK, CONTINUE, GOTO, FALLTHROUGH
		Label string
	}

	DeferStmt struct {
		Call *CallExpr
	}

	GoStmt struct {
		Call *CallExpr
	}

	IncDecStmt struct {
		X  Expr
		Op string // "++" or "--"
	}

	SwitchStmt struct {
		Init Stmt
		Tag  Expr
		Body *BlockStmt
	}

	CaseClause struct {
		List []Expr
		Body []Stmt
	}

	TypeSwitchStmt struct {
		Init   Stmt
		Assign Stmt // e.g., x := y.(type)
		Body   *BlockStmt
	}

	TypeCaseClause struct {
		Types []Type
		Body  []Stmt
	}

	UnsupportedStmt struct {
		Kind string
		Text string
	}
)

func (*BlockStmt) stmtNode()       {}
func (*ExprStmt) stmtNode()        {}
func (*DeclStmt) stmtNode()        {}
func (*AssignStmt) stmtNode()      {}
func (*ReturnStmt) stmtNode()      {}
func (*IfStmt) stmtNode()          {}
func (*ForStmt) stmtNode()         {}
func (*RangeStmt) stmtNode()       {}
func (*SendStmt) stmtNode()        {}
func (*LabeledStmt) stmtNode()     {}
func (*BranchStmt) stmtNode()      {}
func (*DeferStmt) stmtNode()       {}
func (*GoStmt) stmtNode()          {}
func (*IncDecStmt) stmtNode()      {}
func (*SwitchStmt) stmtNode()      {}
func (*CaseClause) stmtNode()      {}
func (*TypeSwitchStmt) stmtNode()  {}
func (*TypeCaseClause) stmtNode()  {}
func (*UnsupportedStmt) stmtNode() {}

func (s *BlockStmt) String() string       { return "{...}" }
func (s *ExprStmt) String() string        { return fmt.Sprintf("%s", s.X) }
func (s *DeclStmt) String() string        { return "decl" }
func (s *AssignStmt) String() string      { return "assign" }
func (s *ReturnStmt) String() string      { return "return" }
func (s *IfStmt) String() string          { return "if" }
func (s *ForStmt) String() string         { return "for" }
func (s *RangeStmt) String() string       { return "range" }
func (s *SendStmt) String() string        { return "send" }
func (s *LabeledStmt) String() string     { return s.Label }
func (s *BranchStmt) String() string      { return s.Tok }
func (s *DeferStmt) String() string       { return "defer" }
func (s *GoStmt) String() string          { return "go" }
func (s *IncDecStmt) String() string      { return "incdec" }
func (s *SwitchStmt) String() string      { return "switch" }
func (s *CaseClause) String() string      { return "case" }
func (s *TypeSwitchStmt) String() string  { return "type switch" }
func (s *TypeCaseClause) String() string  { return "type case" }
func (s *UnsupportedStmt) String() string { return s.Kind }
