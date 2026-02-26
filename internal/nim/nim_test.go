package nim

import (
	"os"
	"os/exec"
	"testing"
)

func TestNimGenerator(t *testing.T) {
	file := &File{
		Nodes: []Node{
			&ProcDecl{
				Name: "hello",
				Args: []Arg{{Name: "name", Typ: "string"}},
				Body: []Node{
					&CallExpr{Fun: "echo", Args: []string{"\"Hello, \" & name"}},
				},
			},
			&CallExpr{Fun: "hello", Args: []string{"\"world\""}},
		},
	}

	code := file.Render(0)
	tmpFile := "test_hello.nim"
	if err := os.WriteFile(tmpFile, []byte(code), 0644); err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpFile)

	cmd := exec.Command("nim", "check", tmpFile)
	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("nim check failed: %v\nOutput:\n%s", err, string(output))
	}
}
