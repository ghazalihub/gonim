package ir

type Symbol struct {
	Name string
	Typ  Type
	Decl Decl
}

type SymbolTable struct {
	Parent  *SymbolTable
	Symbols map[string]*Symbol
}

func NewSymbolTable(parent *SymbolTable) *SymbolTable {
	return &SymbolTable{
		Parent:  parent,
		Symbols: make(map[string]*Symbol),
	}
}

func (s *SymbolTable) Insert(name string, sym *Symbol) {
	s.Symbols[name] = sym
}

func (s *SymbolTable) Lookup(name string) *Symbol {
	if sym, ok := s.Symbols[name]; ok {
		return sym
	}
	if s.Parent != nil {
		return s.Parent.Lookup(name)
	}
	return nil
}
