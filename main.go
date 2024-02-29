package main

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/parser"
	"go/printer"
	"go/token"
	"os"
	"sort"
	"strings"
)

func isStdLib(pkg string) bool {
	// Remove the double quotes from the import path
	pkg = strings.Trim(pkg, "\"")
	// Standard library packages do not contain a dot, while external packages do.
	// This is a heuristic and might not cover every edge case.
	return !strings.Contains(pkg, ".")
}

func main() {
	internalPath := os.Getenv("GO_INTERNAL_PATH") // Use environment variable to define internal imports path
	if internalPath == "" {
		fmt.Println("Environment variable GO_INTERNAL_PATH is not set.")
		os.Exit(1)
	}

	if len(os.Args) < 2 {
		fmt.Println("Usage: sort-imports <filename.go>")
		os.Exit(1)
	}

	filename := os.Args[1]

	// Read the source file
	source, err := os.ReadFile(filename)
	if err != nil {
		fmt.Printf("Error reading file: %s\n", err)
		os.Exit(1)
	}

	// Parse the source file
	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, "", source, parser.ParseComments)
	if err != nil {
		fmt.Printf("Error parsing file: %s\n", err)
		os.Exit(1)
	}

	// Process imports
	for _, decl := range file.Decls {
		genDecl, ok := decl.(*ast.GenDecl)
		if !ok || genDecl.Tok != token.IMPORT {
			continue
		}

		var stdLibImports, externalImports, internalImports []string

		for _, spec := range genDecl.Specs {
			importSpec, ok := spec.(*ast.ImportSpec)
			if !ok {
				continue
			}

			path := importSpec.Path.Value
			if strings.Contains(path, internalPath) {
				internalImports = append(internalImports, path)
			} else if isStdLib(path) {
				stdLibImports = append(stdLibImports, path)
			} else {
				externalImports = append(externalImports, path)
			}
		}

		// Sort the slices
		sort.Strings(stdLibImports)
		sort.Strings(externalImports)
		sort.Strings(internalImports)

		// Combine the sorted slices
		allImports := make([]string, 0, len(stdLibImports)+len(externalImports)+len(internalImports))
		if len(stdLibImports) != 0 {
			allImports = append(allImports, stdLibImports...)
		}
		if len(externalImports) != 0 {
			if len(allImports) > 0 {
				allImports = append(allImports, "")
			}

			allImports = append(allImports, externalImports...)
		}
		if len(internalImports) != 0 {
			if len(allImports) > 0 {
				allImports = append(allImports, "")
			}

			allImports = append(allImports, internalImports...)
		}

		// Construct the new import block
		var newImportBlock bytes.Buffer
		for _, imp := range allImports {
			newImportBlock.WriteString("\t" + imp + "\n")
		}

		// Replace the old import specs with the new sorted ones
		genDecl.Specs = make([]ast.Spec, len(allImports))
		for i, imp := range allImports {
			genDecl.Specs[i] = &ast.ImportSpec{
				Path: &ast.BasicLit{Value: imp},
			}
		}

		// Stop processing after the first import block to avoid duplicating imports
		break
	}

	// Write the modified AST back to a buffer
	var buf bytes.Buffer
	if err := printer.Fprint(&buf, fset, file); err != nil {
		fmt.Printf("Error writing file: %s\n", err)
		os.Exit(1)
	}

	// Write the modified source back to the file
	if err := os.WriteFile(filename, buf.Bytes(), 0644); err != nil {
		fmt.Printf("Error writing file: %s\n", err)
		os.Exit(1)
	}
}
