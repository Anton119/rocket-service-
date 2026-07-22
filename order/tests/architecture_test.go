package tests

import (
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// forbiddenImportRules — правила зависимостей clean architecture для order-модуля.
var forbiddenImportRules = []struct {
	dirPrefix   string
	forbidden   string
	description string
}{
	{
		dirPrefix:   "internal/service",
		forbidden:   "/internal/api/",
		description: "service не должен импортировать api",
	},
	{
		dirPrefix:   "internal/repository",
		forbidden:   "/internal/service/",
		description: "repository не должен импортировать service",
	},
}

func TestArchitecture_DependencyRules(t *testing.T) {
	t.Helper()

	moduleRoot, err := findModuleRoot()
	if err != nil {
		t.Fatalf("найти корень модуля order: %v", err)
	}

	fset := token.NewFileSet()
	var violations []string

	for _, rule := range forbiddenImportRules {
		root := filepath.Join(moduleRoot, rule.dirPrefix)
		walkErr := filepath.WalkDir(root, func(path string, d os.DirEntry, walkErr error) error {
			if walkErr != nil {
				return walkErr
			}
			if d.IsDir() {
				if d.Name() == "mocks" || d.Name() == "tests" {
					return filepath.SkipDir
				}
				return nil
			}
			if !strings.HasSuffix(path, ".go") || strings.HasSuffix(path, "_test.go") {
				return nil
			}

			file, parseErr := parser.ParseFile(fset, path, nil, parser.ImportsOnly)
			if parseErr != nil {
				return parseErr
			}

			for _, imp := range file.Imports {
				importPath := strings.Trim(imp.Path.Value, `"`)
				if strings.Contains(importPath, rule.forbidden) {
					violations = append(violations,
						filepath.Join(rule.dirPrefix, filepath.Base(path))+": "+rule.description+
							" (import "+importPath+")",
					)
				}
			}
			return nil
		})
		if walkErr != nil {
			t.Fatalf("обход %s: %v", rule.dirPrefix, walkErr)
		}
	}

	if len(violations) > 0 {
		t.Fatalf("нарушения архитектурных правил:\n- %s", strings.Join(violations, "\n- "))
	}
}

func findModuleRoot() (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", err
	}

	for {
		if _, statErr := os.Stat(filepath.Join(dir, "go.mod")); statErr == nil {
			return dir, nil
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			return "", os.ErrNotExist
		}
		dir = parent
	}
}
