package main

import (
	"go/parser"
	"go/token"
	"os"
	"os/exec"
	"strings"

	"github.com/clipperhouse/typewriter"
)

// get runs `go get` for required typewriters, either default or specified in _gen.go
func get(args []string) error {
	imports, err := getTypewriterImports()

	if err != nil {
		return err
	}

	// we just want the paths
	imps := imports.SelectString(func(imp typewriter.ImportSpec) string {
		return imp.Path
	})

	get := []string{"get"}
	get = append(get, args...)
	get = append(get, imps...)

	cmd := exec.Command("go", get...)
	cmd.Stdout = out
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return nil
	}

	return nil
}

func getTypewriterImports() (typewriter.ImportSpecSlice, error) {
	var imports typewriter.ImportSpecSlice

	// check for existence of custom file
	if src, err := os.Open(customName); err == nil {
		defer src.Close()

		// custom file exists, parse its imports
		fset := token.NewFileSet()
		f, err := parser.ParseFile(fset, "", src, parser.ImportsOnly)
		if err != nil {
			return imports, err
		}
		for _, v := range f.Imports {
			imp := typewriter.ImportSpec{
				Name: v.Name.Name,
				Path: strings.Trim(v.Path.Value, `"`), // lose the quotes
			}
			imports = append(imports, imp)
		}
	} else {
		// doesn't exist, use standard
		imports = append(imports, stdImports...)
	}

	return imports, nil
}
