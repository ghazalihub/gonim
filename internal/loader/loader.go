package loader

import (
	"fmt"

	"github.com/user/go2nim/internal/ir"
	"golang.org/x/tools/go/packages"
)

func Load(patterns ...string) (*ir.Program, error) {
	cfg := &packages.Config{
		Mode: packages.NeedName | packages.NeedFiles | packages.NeedCompiledGoFiles |
			packages.NeedImports | packages.NeedTypes | packages.NeedTypesInfo |
			packages.NeedSyntax | packages.NeedTypesSizes,
	}
	pkgs, err := packages.Load(cfg, patterns...)
	if err != nil {
		return nil, err
	}

	if packages.PrintErrors(pkgs) > 0 {
		return nil, fmt.Errorf("errors loading packages")
	}

	program := &ir.Program{}
	for _, pkg := range pkgs {
		irPkg, err := convertPackage(pkg)
		if err != nil {
			return nil, err
		}
		program.Packages = append(program.Packages, irPkg)
	}

	return program, nil
}
