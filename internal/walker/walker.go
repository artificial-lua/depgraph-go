package walker

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/artificial-lua/depgraph-go/internal/parser"
)

// CollectPackages walks through the project and collects package dependencies
func CollectPackages(rootDir, moduleName string) (map[string][]string, error) {
	packageRefs := make(map[string][]string)

	err := filepath.Walk(rootDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			if strings.HasSuffix(info.Name(), "_test") {
				return filepath.SkipDir
			}
			return nil
		}

		if !strings.HasSuffix(info.Name(), ".go") || strings.HasSuffix(info.Name(), "_test.go") {
			return nil
		}

		dir := filepath.Dir(path)
		imports, err := parser.ParseGoFile(path)
		if err != nil {
			return err
		}

		internalImports := parser.StripModulePrefix(imports, moduleName)
		if len(internalImports) > 0 {
			relPath, _ := filepath.Rel(rootDir, dir)
			if relPath == "." {
				relPath = ""
			}
			if strings.Contains(relPath, "_test") {
				return nil
			}

			existing := packageRefs[relPath]
			unique := make(map[string]bool)
			for _, v := range existing {
				unique[v] = true
			}
			for _, v := range internalImports {
				if !unique[v] {
					existing = append(existing, v)
					unique[v] = true
				}
			}
			packageRefs[relPath] = existing
		}
		return nil
	})

	return packageRefs, err
}
