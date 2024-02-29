package main

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"os"
	"sort"
	"strconv"
	"strings"
)

const newLinePlaceholder = "// NEW LINE PLACEHOLDER (DELETE)"

// Import represents a single import statement, including its alias (if any).
type Import struct {
	Alias string
	Path  string
}

// ImportsByType categorizes imports into standard, external, and internal.
type ImportsByType struct {
	Standard []Import
	External []Import
	Internal []Import
}

// isStdLib checks whether the import is from golang stdlib or not.
func isStdLib(importPath string) bool {
	return !strings.Contains(importPath, ".")
}

// isInternal checks whether the import is internal for the specific project.
func isInternal(importPath, internalPath string) bool {
	return strings.HasPrefix(importPath, strings.Trim(internalPath, "\""))
}

// addSortedImportsToAST adds the sorted imports to abstract syntax tree.
func addSortedImportsToAST(node *ast.File, sortedImports ImportsByType) {
	var specs []ast.Spec

	// Helper function to append imports to specs.
	appendImports := func(imports []Import) {
		for _, imp := range imports {
			quotedPath := strconv.Quote(imp.Path)
			spec := &ast.ImportSpec{
				Path: &ast.BasicLit{
					Kind:  token.STRING,
					Value: quotedPath,
				},
			}
			if imp.Alias != "" {
				spec.Name = ast.NewIdent(imp.Alias)
			}
			specs = append(specs, spec)
		}
	}

	// Append standard imports.
	appendImports(sortedImports.Standard)

	// If there are standard imports and either external or internal imports, add a newline separator.
	if len(sortedImports.Standard) > 0 && (len(sortedImports.External) > 0 || len(sortedImports.Internal) > 0) {
		specs = append(specs, &ast.ImportSpec{Path: &ast.BasicLit{Kind: token.COMMENT, Value: newLinePlaceholder}})
	}

	// Append external imports.
	appendImports(sortedImports.External)

	// If there are external imports and internal imports, add a newline separator.
	if len(sortedImports.External) > 0 && len(sortedImports.Internal) > 0 {
		specs = append(specs, &ast.ImportSpec{Path: &ast.BasicLit{Kind: token.COMMENT, Value: newLinePlaceholder}})
	}

	// Append internal imports.
	appendImports(sortedImports.Internal)

	// Only create an import declaration if there are any imports to add.
	if len(specs) > 0 {
		importDecl := &ast.GenDecl{
			Tok:   token.IMPORT,
			Specs: specs,
		}
		if len(specs) > 1 {
			importDecl.Lparen = 1 // Indicate that imports are in a block if there's more than one import.
		}

		// Prepend the combined import declaration to the file's declaration list.
		node.Decls = append([]ast.Decl{importDecl}, node.Decls...)
	}
}

func sortImports(imports []Import) {
	sort.Slice(imports, func(i, j int) bool {
		return imports[i].Path < imports[j].Path
	})
}

func main() {
	internalPath := os.Getenv("GO_INTERNAL_PATH") // Use environment variable to define internal imports path
	if internalPath == "" {
		fmt.Println("Environment variable GO_INTERNAL_PATH is not set.")
		os.Exit(1)
	}

	if len(os.Args) != 2 {
		fmt.Println("Usage: sortimports <filename>")
		os.Exit(1)
	}
	filename := os.Args[1]

	// Parse the Go file to construct an AST.
	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, filename, nil, parser.ParseComments)
	if err != nil {
		fmt.Printf("Error parsing file: %s\n", err)
		os.Exit(1)
	}

	// Initialize a structure to hold categorized imports.
	var sortedImports ImportsByType

	// Collect and categorize all imports from the AST.
	ast.Inspect(node, func(n ast.Node) bool {
		// Find import declarations.
		if genDecl, ok := n.(*ast.GenDecl); ok && genDecl.Tok == token.IMPORT {
			for _, spec := range genDecl.Specs {
				importSpec, ok := spec.(*ast.ImportSpec)
				if !ok {
					continue
				}
				imp := Import{
					Path:  strings.Trim(importSpec.Path.Value, "\""),
					Alias: "",
				}
				if importSpec.Name != nil { // Capture the alias if present.
					imp.Alias = importSpec.Name.Name
				}

				// Categorize the import based on its path.
				if isStdLib(imp.Path) {
					sortedImports.Standard = append(sortedImports.Standard, imp)
				} else if isInternal(imp.Path, internalPath) {
					sortedImports.Internal = append(sortedImports.Internal, imp)
				} else {
					sortedImports.External = append(sortedImports.External, imp)
				}
			}
			return false // No need to inspect deeper since we're only interested in top-level imports.
		}
		return true
	})

	sortImports(sortedImports.Standard)
	sortImports(sortedImports.External)
	sortImports(sortedImports.Internal)

	// Remove the old import declarations from the AST.
	newDecls := make([]ast.Decl, 0)
	for _, decl := range node.Decls {
		if genDecl, ok := decl.(*ast.GenDecl); ok && genDecl.Tok == token.IMPORT {
			continue
		}
		newDecls = append(newDecls, decl)
	}
	node.Decls = newDecls

	// Write back the sorted and categorized imports into the AST as a single import block.
	addSortedImportsToAST(node, sortedImports)

	// Write the modified AST back to the file.
	var buf bytes.Buffer
	if err := format.Node(&buf, fset, node); err != nil {
		fmt.Printf("Error formatting node: %s\n", err)
		os.Exit(1)
	}

	formattedCode := buf.String()
	// Replace the newline placeholders with actual newlines.
	formattedCodeWithNewlines := strings.Replace(formattedCode, newLinePlaceholder, "\n", -1)

	if err := os.WriteFile(filename, []byte(formattedCodeWithNewlines), 0644); err != nil {
		fmt.Printf("Error writing file: %s\n", err)
		os.Exit(1)
	}
}
