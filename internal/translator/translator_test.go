package translator

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/user/go2nim/internal/loader"
)

func TestBasicTranspilation(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "go2nim-transpilation-test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	goFile := filepath.Join(tmpDir, "main.go")
	content := `package main
func add(a, b int) int {
	return a + b
}
func main() {
	var x = add(1, 2)
	if x == 3 {
		println("ok")
	}
}
`
	if err := os.WriteFile(goFile, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	prog, err := loader.Load(goFile)
	if err != nil {
		t.Fatalf("Load failed: %v", err)
	}

	nimFiles := Translate(prog)
	if len(nimFiles) == 0 {
		t.Fatal("no nim files generated")
	}

	nimCode := nimFiles[0].Render(0)
	t.Logf("Generated Nim code:\n%s", nimCode)

	// Optionally try to compile with nim check
}
