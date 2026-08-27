package translator

import (
	"os"
	"path/filepath"
	"strings"
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

func TestStressLanguageConstructs(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "go2nim-stress-test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	goFile := filepath.Join(tmpDir, "main.go")
	content := `package main

type pair struct { a, b int }

var x, y = 1, 2
const a, b = 3, 4

func tuple() (int, int) { return x, y }

func main() {
	m := map[string]int{"a": 1}
	delete(m, "a")
	s := []int{1, 2, 3}
	_ = s[1:]
	for i := 0; i < 3; i++ { x += i }
	for k, v := range m { println(k, v) }
	if z := x + y; z > 0 { println(z) }
	switch {
	case x > y:
		println("x")
	default:
		println("y")
	}
	x, y = tuple()
}
`
	if err := os.WriteFile(goFile, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	prog, err := loader.Load(goFile)
	if err != nil {
		t.Fatalf("Load failed: %v", err)
	}
	nimCode := Translate(prog)[0].Render(0)
	checks := []string{
		"var x: int = 1",
		"var y: int = 2",
		"const a = 3",
		"const b = 4",
		"m.del(\"a\")",
		"s[1..^1]",
		"case true",
		"(x, y) = tupleX()",
	}
	for _, want := range checks {
		if !strings.Contains(nimCode, want) {
			t.Fatalf("generated Nim missing %q:\n%s", want, nimCode)
		}
	}
}
