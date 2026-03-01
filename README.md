# go2nim: Production-Grade Go → Nim Transpiler

`go2nim` is a compiler-grade transpilation tool designed to convert complete Go packages into valid, compilable Nim code while preserving Go's semantics and type safety.

## 🏗 Architecture

`go2nim` follows a standard compiler pipeline:
1. **Frontend**: Uses `go/packages`, `go/ast`, and `go/types` to load, parse, and type-check Go source.
2. **IR Generation**: Converts Typed Go AST into a custom, language-neutral **Intermediate Representation (IR)**.
3. **Transformation**: Executes passes to translate Go-specific constructs into Nim equivalents.
4. **Nim AST**: Construct a Nim-specific Abstract Syntax Tree.
5. **Backend**: A code generator that renders Nim source code with proper indentation and syntax.

## 🚀 Features Checklist

### ✅ Supported Core Features
- [x] **Multi-file Packages**: Correctly resolves symbols and types across multiple files in a package.
- [x] **Local Imports**: Maps local Go imports to Nim modules.
- [x] **Type System**:
    - [x] Basic types (int, string, bool, etc.)
    - [x] Named types and Type aliases.
    - [x] Structs (with field tags and exported/unexported fields).
    - [x] Struct Embedding (Anonymous fields).
    - [x] Interfaces (mapped to Nim `concepts`).
    - [x] Pointers, Slices, and Arrays.
    - [x] Maps (mapped to Nim `Table`, with automatic `import tables`).
- [x] **Methods & Receivers**: Supports both value and pointer receivers.
- [x] **Generics**: Fully supports type parameters and constraints (Go 1.18+).
- [x] **Constants**:
    - [x] Typed and Untyped constants.
    - [x] `iota` resolution.
- [x] **Functions**:
    - [x] Variadic parameters (`...T` → `varargs[T]`).
    - [x] Multiple return values (mapped to Nim tuples).
    - [x] Named return parameters (mapped to Nim's `result`).
- [x] **Control Flow**:
    - [x] `if-else` (including initialization statements).
    - [x] `for` loops (standard and infinite).
    - [x] `for range` (over slices, arrays, and maps).
    - [x] `switch` and `type switch`.
    - [x] Labels and `break/continue` to labels.

### ✅ Advanced Go Features
- [x] **Defer**: Mapped to Nim's `defer:` block.
- [x] **Panic/Recover**: Basic mapping to Nim exceptions.
- [x] **cgo**: Captures C preamble and maps `C.xxx` calls using Nim's `{.emit.}`.
- [x] **go:embed**: Maps embedded assets to Nim's `staticRead`.
- [x] **Blank Identifier**: Correctly handles `_` in assignments and imports.

### ✅ Builtin Functions
- [x] `len`, `cap`
- [x] `make`, `new`
- [x] `append`, `copy`, `delete`
- [x] `println`, `print` (mapped to `echo`/`write`)

## ❌ Not Supported (Yet)
- [ ] **Goroutines**: No current mapping for `go` statements to Nim threads.
- [ ] **Channels**: No mapping for Go channels or `select` blocks.
- [ ] **Standard Library Rewriting**: Currently relies on dummy stubs in `gostdnim/`.
- [ ] **Complex/Real/Imag**: Basic mapping missing.

## 🛠 Usage

To transpile a Go package:

```bash
go run cmd/go2nim/main.go -o ./output_dir ./path/to/your/go/package
```

### 📦 Dummy Stdlib
The transpiler includes a `gostdnim/` folder containing dummy Nim modules for common Go packages (`fmt`, `os`, `strings`, etc.). To compile the generated Nim code, ensure `gostdnim` is in your Nim import path.

## 🗺 Roadmap / TODO
- [ ] Implement Go-style concurrency (Goroutines and Channels) using Nim's `asyncdispatch` or `threading`.
- [ ] Expand `gostdnim` to include real logic for core Go standard library functions.
- [ ] Support for `complex64` and `complex128` types.
- [ ] Optimization passes for Nim-specific idioms (e.g., using `let` where appropriate).
- [ ] Improved error reporting and source mapping.

## 🧠 Design Philosophy
- **Accuracy First**: Rely on Go's official type-checker to resolve all symbols before translation.
- **Independence**: The IR layer ensures that we aren't tied to Go's internal AST structure, allowing for easier future transformations.
- **Readability**: Generate Nim code that is as close to the original Go as possible while remaining valid Nim.
