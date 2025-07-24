package parser

import (
	"bufio"
	"os"
	"regexp"
	"strings"
)

// GetModuleName reads go.mod and extracts the module name
func GetModuleName(goModPath string) (string, error) {
	file, err := os.Open(goModPath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if strings.HasPrefix(line, "module ") {
			return strings.TrimSpace(strings.TrimPrefix(line, "module ")), nil
		}
	}
	return "", nil
}

// ParseGoFile extracts imported packages from a .go file
func ParseGoFile(path string) ([]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var imports []string
	importBlock := false
	re := regexp.MustCompile(`"([^"]+)"`)

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if strings.HasPrefix(line, "import (") {
			importBlock = true
		} else if importBlock {
			if strings.HasPrefix(line, ")") {
				importBlock = false
			} else {
				matches := re.FindStringSubmatch(line)
				if len(matches) > 1 {
					imports = append(imports, matches[1])
				}
			}
		} else if strings.HasPrefix(line, "import ") {
			matches := re.FindStringSubmatch(line)
			if len(matches) > 1 {
				imports = append(imports, matches[1])
			}
		}
	}
	return imports, nil
}

// StripModulePrefix removes the module prefix from imports
func StripModulePrefix(imports []string, moduleName string) []string {
	var stripped []string
	for _, imp := range imports {
		if strings.HasPrefix(imp, moduleName) {
			path := strings.TrimPrefix(imp, moduleName)
			path = strings.TrimPrefix(path, "/")
			stripped = append(stripped, path)
		}
	}
	return stripped
}
