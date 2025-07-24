package main

import (
	"fmt"
	"os"

	"github.com/artificial-lua/depgraph-go/internal/graph" // TODO: Reverse arrow direction in GenerateDot implementation
	"github.com/artificial-lua/depgraph-go/internal/parser"
	"github.com/artificial-lua/depgraph-go/internal/walker"
)

func main() {
	goModPath := "go.mod"
	if _, err := os.Stat(goModPath); os.IsNotExist(err) {
		fmt.Println("go.mod file not found. Run this program at the Go project root.")
		return
	}

	moduleName, err := parser.GetModuleName(goModPath)
	if err != nil {
		fmt.Println("Error reading module name:", err)
		return
	}

	packages, err := walker.CollectPackages(".", moduleName)
	if err != nil {
		fmt.Println("Error collecting packages:", err)
		return
	}

	if err := graph.GenerateDot(packages, "graph.dot"); err != nil {
		fmt.Println("Error generating dot file:", err)
		return
	}

	fmt.Println("Graphviz .dot file generated as graph.dot")
}
