package loader

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/user/go2nim/internal/ir"
)

func TestLoadSimple(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "go2nim-test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	goFile := filepath.Join(tmpDir, "main.go")
	content := `package main
func main() {
	var x int = 42
	println(x)
}
`
	if err := os.WriteFile(goFile, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	// We need a go.mod in the tmpDir if we are in module mode,
	// or just run in the same module if we refer to it.
	// For simplicity, we'll just try to load the file.

	prog, err := Load(goFile)
	if err != nil {
		t.Fatalf("Load failed: %v", err)
	}

	if len(prog.Packages) == 0 {
		t.Fatal("no packages loaded")
	}

	pkg := prog.Packages[0]
	if pkg.Name != "main" {
		t.Errorf("expected package main, got %s", pkg.Name)
	}

	if len(pkg.Files) == 0 {
		t.Fatal("no files in package")
	}

	foundMain := false
	for _, decl := range pkg.Files[0].Decls {
		if f, ok := decl.(*ir.FuncDecl); ok && f.Name == "main" {
			foundMain = true
			if len(f.Body.List) != 2 {
				t.Errorf("expected 2 statements in main, got %d", len(f.Body.List))
			}
		}
	}

	if !foundMain {
		t.Error("func main not found in IR")
	}
}
