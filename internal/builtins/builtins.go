package builtins

type BuiltinInfo struct {
	NimName string
	IsMember bool
}

var Builtins = map[string]BuiltinInfo{
	"len":     {NimName: "len", IsMember: false},
	"cap":     {NimName: "cap", IsMember: false},
	"append":  {NimName: "add", IsMember: true},
	"make":    {NimName: "newSeq", IsMember: false},
	"new":     {NimName: "new", IsMember: false},
	"delete":  {NimName: "del", IsMember: false},
	"copy":    {NimName: "copy", IsMember: false},
	"panic":   {NimName: "raise", IsMember: false},
	"println": {NimName: "echo", IsMember: false},
	"print":   {NimName: "write", IsMember: false},
}
