package types

import (
	"fmt"
	"strings"

	"github.com/user/go2nim/internal/ir"
)

type ImportAdder interface {
	AddImport(name string)
}

func MapType(t ir.Type, adder ImportAdder) string {
	switch tt := t.(type) {
	case *ir.BasicType:
		switch tt.Name {
		case "int":
			return "int"
		case "int32":
			return "int32"
		case "int64":
			return "int64"
		case "uint":
			return "uint"
		case "uint32":
			return "uint32"
		case "uint64":
			return "uint64"
		case "float32":
			return "float32"
		case "float64":
			return "float64"
		case "string":
			return "string"
		case "bool":
			return "bool"
		case "byte":
			return "byte"
		case "rune":
			return "int32"
		case "uintptr":
			return "uintptr"
		default:
			return tt.Name
		}
	case *ir.PointerType:
		return fmt.Sprintf("ref %s", MapType(tt.Elem, adder))
	case *ir.SliceType:
		return fmt.Sprintf("seq[%s]", MapType(tt.Elem, adder))
	case *ir.ArrayType:
		return fmt.Sprintf("array[%d, %s]", tt.Len, MapType(tt.Elem, adder))
	case *ir.MapType:
		if adder != nil {
			adder.AddImport("tables")
		}
		return fmt.Sprintf("Table[%s, %s]", MapType(tt.Key, adder), MapType(tt.Value, adder))
	case *ir.NamedType:
		if len(tt.TypeArgs) > 0 {
			var args []string
			for _, arg := range tt.TypeArgs {
				args = append(args, MapType(arg, adder))
			}
			return fmt.Sprintf("%s[%s]", tt.Name, strings.Join(args, ", "))
		}
		return tt.Name
	case *ir.StructType:
		return "object" // Simplified
	case *ir.InterfaceType:
		return "concept" // Simplified
	case *ir.FuncType:
		return "proc" // Simplified
	case *ir.TypeParam:
		return tt.Name
	}
	return "auto"
}
