package graph

import (
	"fmt"
	"os"
)

// GenerateDot creates a Graphviz .dot file
func GenerateDot(packages map[string][]string, outputFile string) error {
	file, err := os.Create(outputFile)
	if err != nil {
		return err
	}
	defer file.Close()

	file.WriteString("digraph G {\n")
	file.WriteString("    rankdir=LR;\n")
	file.WriteString("    node [shape=box, style=filled, fillcolor=lightblue];\n")

	for pkg, refs := range packages {
		pkgName := pkg
		if pkgName == "" {
			pkgName = "."
		}
		for _, ref := range refs {
			file.WriteString(fmt.Sprintf("    \"%s\" -> \"%s\";\n", pkgName, ref))
		}
	}

	file.WriteString("}\n")
	return nil
}
