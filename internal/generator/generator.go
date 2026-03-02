package generator

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/DiegoDev2/basalt/pkg/ast"
)

type Generator struct {
	file       *ast.File
	outDir     string
	OnProgress func(string)
}

func NewGenerator(file *ast.File, outDir string) *Generator {
	return &Generator{file: file, outDir: outDir}
}

func (g *Generator) log(msg string) {
	if g.OnProgress != nil {
		g.OnProgress(msg)
	}
}

func (g *Generator) Generate() error {
	g.log("Parsing main.bs...")

	if err := os.MkdirAll(g.outDir, 0755); err != nil {
		return err
	}

	g.log(fmt.Sprintf("AST built — %d tables, %d resources", len(g.file.Tables), len(g.file.Resources)))
	g.log("\nGenerating shared")

	sharedFiles := []string{
		"lib/supabase.ts",
		"lib/auth.ts",
		"middleware/error.ts",
		"middleware/logger.ts",
		"index.ts",
		"schema.sql",
		"package.json",
		"tsconfig.json",
		".env.example",
	}

	for _, relPath := range sharedFiles {
		if err := g.generateSharedFile(relPath); err != nil {
			return err
		}
		g.log(fmt.Sprintf(" ✓ %-20s", relPath))
	}

	for _, res := range g.file.Resources {
		g.log("\nGenerating " + res.Name)

		resDir := filepath.Join(g.outDir, res.Name)
		if err := os.MkdirAll(resDir, 0755); err != nil {
			return err
		}

		resourceFiles := []string{"router.ts", "handlers.ts", "model.ts", "supabase.ts"}
		for _, relPath := range resourceFiles {
			if err := g.generateResourceFile(res, relPath); err != nil {
				return err
			}
			g.log(fmt.Sprintf(" ✓ %-20s", relPath))
		}
	}

	return nil
}

func (g *Generator) generateSharedFile(relPath string) error {
	fullPath := filepath.Join(g.outDir, relPath)

	if err := os.MkdirAll(filepath.Dir(fullPath), 0755); err != nil {
		return err
	}

	tmplStr, ok := SharedTemplates[relPath]
	if !ok {
		return fmt.Errorf("template not found: %s", relPath)
	}

	tmpl, err := template.New(relPath).Funcs(template.FuncMap{
		"lower":      strings.ToLower,
		"capitalize": capitalize,
		"replace":    strings.ReplaceAll,
	}).Parse(tmplStr)
	if err != nil {
		return err
	}

	f, err := os.Create(fullPath)
	if err != nil {
		return err
	}
	defer f.Close()

	return tmpl.Execute(f, g.file)
}

func (g *Generator) generateResourceFile(res *ast.Resource, relPath string) error {
	fullPath := filepath.Join(g.outDir, res.Name, relPath)

	tmplStr, ok := ResourceTemplates[relPath]
	if !ok {
		return fmt.Errorf("resource template not found: %s", relPath)
	}

	tmpl, err := template.New(relPath).Funcs(template.FuncMap{
		"lower":      strings.ToLower,
		"capitalize": capitalize,
		"replace":    strings.ReplaceAll,
		// trimResourcePrefix strips the resource name from the path.
		// e.g. "/users/:id" -> "/:id", "/users" -> "/"
		"trimResourcePrefix": func(path, resourceName string) string {
			prefix := "/" + resourceName
			if path == prefix {
				return "/"
			}
			return strings.TrimPrefix(path, prefix)
		},
		// handlerName generates a valid TypeScript identifier from method + path.
		// e.g. ("GET",    "/users",    "users") -> "getAll"
		// e.g. ("POST",   "/users",    "users") -> "post"
		// e.g. ("GET",    "/users/:id","users") -> "get_id"
		// e.g. ("PUT",    "/users/:id","users") -> "put_id"
		// e.g. ("DELETE", "/users/:id","users") -> "delete_id"
		"handlerName": func(method, path, resourceName string) string {
			prefix := "/" + resourceName
			stripped := strings.TrimPrefix(path, prefix)

			if stripped == "" || stripped == "/" {
				if strings.ToUpper(method) == "GET" {
					return "getAll"
				}
				return strings.ToLower(method)
			}

			// Remove slashes and colons to produce a clean identifier
			name := strings.ReplaceAll(stripped, "/", "_")
			name = strings.ReplaceAll(name, ":", "")
			name = strings.Trim(name, "_")

			return strings.ToLower(method) + "_" + name
		},
	}).Parse(tmplStr)
	if err != nil {
		return err
	}

	f, err := os.Create(fullPath)
	if err != nil {
		return err
	}
	defer f.Close()

	data := struct {
		Resource *ast.Resource
		File     *ast.File
		Config   *ast.Config
		Table    *ast.Table
	}{
		Resource: res,
		File:     g.file,
		Config:   g.file.Config,
	}

	for _, t := range g.file.Tables {
		if t.Name == res.Table {
			data.Table = t
			break
		}
	}

	return tmpl.Execute(f, data)
}

func capitalize(s string) string {
	if len(s) == 0 {
		return s
	}
	return strings.ToUpper(s[0:1]) + s[1:]
}
