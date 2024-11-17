package internal

import (
	"bufio"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

// GenerateMermaidFromFile generates a Mermaid state diagram from a given Go source file.
func GenerateMermaidFromFile(filePath string) (string, error) {
	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, filePath, nil, parser.AllErrors)
	if err != nil {
		return "", fmt.Errorf("failed to parse file %s: %w", filePath, err)
	}

	var chart strings.Builder
	chart.WriteString("stateDiagram-v2\n")

	commentRegex := regexp.MustCompile(`^//\s*([\w, ]+)\s*->\s*([\w]+)$`)

	for _, decl := range file.Decls {
		if genDecl, ok := decl.(*ast.GenDecl); ok {
			for _, spec := range genDecl.Specs {
				if typeSpec, ok := spec.(*ast.TypeSpec); ok {
					if interfaceType, ok := typeSpec.Type.(*ast.InterfaceType); ok {
						for _, method := range interfaceType.Methods.List {
							if len(method.Names) > 0 {
								methodName := method.Names[0].Name
								comment := ""
								if method.Doc != nil {
									comment = strings.TrimSpace(method.Doc.Text())
								} else {
									comment = extractCommentFromFile(filePath, method.Pos())
								}

								if matches := commentRegex.FindStringSubmatch(comment); matches != nil {
									fromStates := strings.Split(matches[1], ",")
									toState := matches[2]
									for _, fromState := range fromStates {
										fromState = strings.TrimSpace(fromState)
										chart.WriteString(fmt.Sprintf("    %s -->|%s| %s\n", fromState, methodName, toState))
									}
								}
							}
						}
					}
				}
			}
		}
	}

	return chart.String(), nil
}

// GenerateMermaidFromDirectory scans a directory for _state_machine.go files and generates Mermaid diagrams.
func GenerateMermaidFromDirectory(directory string) (map[string]string, error) {
	results := make(map[string]string)
	err := filepath.Walk(directory, func(path string, info os.FileInfo, err error) error {
		if strings.HasSuffix(path, "_state_machine.go") {
			diagram, err := GenerateMermaidFromFile(path)
			if err != nil {
				return err
			}
			results[path] = diagram
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return results, nil
}

// extractCommentFromFile reads the file and retrieves the comment preceding a method.
func extractCommentFromFile(filePath string, pos token.Pos) string {
	f, err := os.Open(filePath)
	if err != nil {
		return ""
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	lineNumber := 1
	targetLine := positionToLineNumber(filePath, pos)

	for scanner.Scan() {
		if lineNumber == targetLine-1 {
			return strings.TrimSpace(scanner.Text())
		}
		lineNumber++
	}
	return ""
}

// positionToLineNumber converts a token.Pos to a line number in the source file.
func positionToLineNumber(filePath string, pos token.Pos) int {
	fset := token.NewFileSet()
	// _, _ := parser.ParseFile(fset, filePath, nil, parser.AllErrors)
	return fset.Position(pos).Line
}
