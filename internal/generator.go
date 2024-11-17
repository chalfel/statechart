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
					fmt.Printf("Found type: %s\n", typeSpec.Name.Name)

					if interfaceType, ok := typeSpec.Type.(*ast.InterfaceType); ok {
						fmt.Printf("Found interface: %s\n", typeSpec.Name.Name)

						for _, method := range interfaceType.Methods.List {
							if len(method.Names) > 0 {
								methodName := method.Names[0].Name
								comment := ""
								if method.Doc != nil {
									comment = strings.TrimSpace(method.Doc.Text())
								} else {
									comment = extractCommentFromFile(fset, filePath, method.Pos())
								}

								fmt.Printf("Method: %s, Comment: %s\n", methodName, comment)

								if matches := commentRegex.FindStringSubmatch(comment); matches != nil {
									fmt.Printf("Regex Match: %v\n", matches)
									fromStates := strings.Split(matches[1], ",")
									toState := matches[2]
									for _, fromState := range fromStates {
										fromState = strings.TrimSpace(fromState)
										chart.WriteString(fmt.Sprintf("    %s -->|%s| %s\n", fromState, methodName, toState))
									}
								} else {
									fmt.Printf("No match for comment: %s\n", comment)
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

// positionToLineNumber converts a token.Pos to a line number in the source file.
func positionToLineNumber(fset *token.FileSet, pos token.Pos) int {
	return fset.Position(pos).Line
}

// extractCommentFromFile reads the file and retrieves the comment preceding a method.
func extractCommentFromFile(fset *token.FileSet, filePath string, pos token.Pos) string {
	f, err := os.Open(filePath)
	if err != nil {
		return ""
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	targetLine := positionToLineNumber(fset, pos) - 1 // Line preceding the method position
	currentLine := 0

	for scanner.Scan() {
		line := scanner.Text()
		if currentLine == targetLine {
			// Find the `//` and return everything after it
			if idx := strings.Index(line, "//"); idx != -1 {
				return strings.TrimSpace(line[idx:])
			}
			return ""
		}
		currentLine++
	}
	return ""
}

// GenerateMermaidFromInterfaces scans a file for XStateMachine interfaces and generates Mermaid diagrams.
func GenerateMermaidFromInterfaces(filePath string) (map[string]string, error) {
	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, filePath, nil, parser.AllErrors)
	if err != nil {
		return nil, fmt.Errorf("failed to parse file %s: %w", filePath, err)
	}

	diagrams := make(map[string]string)
	commentRegex := regexp.MustCompile(`^//\s*([\w, ]+)\s*->\s*([\w]+)$`)

	for _, decl := range file.Decls {
		if genDecl, ok := decl.(*ast.GenDecl); ok {
			for _, spec := range genDecl.Specs {
				if typeSpec, ok := spec.(*ast.TypeSpec); ok {
					interfaceName := typeSpec.Name.Name
					if strings.HasSuffix(interfaceName, "StateMachine") {
						chart := &strings.Builder{}
						chart.WriteString("stateDiagram-v2\n")
						fmt.Printf("Found interface: %s\n", interfaceName)

						if interfaceType, ok := typeSpec.Type.(*ast.InterfaceType); ok {
							for _, method := range interfaceType.Methods.List {
								if len(method.Names) > 0 {
									methodName := method.Names[0].Name
									comment := ""
									if method.Doc != nil {
										comment = strings.TrimSpace(method.Doc.Text())
									} else {
										comment = extractCommentFromFile(fset, filePath, method.Pos())
									}
									fmt.Printf("Method: %s, Comment: %s\n", methodName, comment)

									if matches := commentRegex.FindStringSubmatch(comment); matches != nil {
										fromStates := strings.Split(matches[1], ",")
										toState := matches[2]
										for _, fromState := range fromStates {
											fromState = strings.TrimSpace(fromState)
											chart.WriteString(fmt.Sprintf("    %s --> %s : %s()\n", fromState, toState, methodName))

										}
									}
								}
							}
						}

						diagrams[interfaceName] = chart.String()
					}
				}
			}
		}
	}

	return diagrams, nil
}
