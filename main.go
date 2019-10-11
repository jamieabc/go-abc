package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"math"
	"os"
	"path/filepath"
)

type report struct {
	path        string
	assignments int
	branches    int
	conditions  int
	score       int
}

func (r report) String() string {
	return fmt.Sprintf(
		"file: %s\n\tassignment: %d\n\tbranch: %d\n\tcondition: %d\n\tscore: %d\n",
		r.path,
		r.assignments,
		r.branches,
		r.conditions,
		r.score,
	)
}

func main() {
	if 1 == len(os.Args) {
		usage()
		return
	}

	source := os.Args[1]
	files, err := listFile(source)
	if nil != err {
		fmt.Println(err)
		return
	}

	var results []report
	for _, f := range files {
		results = append(results, analyze(f))
	}
	fmt.Printf("%s\n", results)
}

func usage() {
	fmt.Println("Usage: go-abc source")
	fmt.Println("source could be file or directory, e.g. main.go, /home/user/test.go, etc.")
}

func listFile(source string) (files []string, err error) {
	info, e := os.Stat(source)
	if nil != e {
		err = fmt.Errorf("open %s with error: %s", source, err)
		return
	}

	root, e := filepath.Abs(source)
	if nil != e {
		err = fmt.Errorf("get absolute root of %s with error: %s", source, e)
		return
	}

	if info.IsDir() {
		err = filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
			if nil != err {
				return fmt.Errorf("walk through dir with error: %s", err)
			}

			if !info.IsDir() {
				match, err := filepath.Match("*.go", filepath.Base(path))
				if !match || err != nil {
					return nil
				}

				files = append(files, path)
			}
			return nil
		})

		if nil != err {
			return
		}
	} else {
		match, e := filepath.Match("*.go", filepath.Base(root))
		if !match || e != nil {
			return
		}
		files = append(files, root)
	}

	return
}

func analyze(file string) (r report) {
	match, err := filepath.Match("*_test.go", filepath.Base(file))
	if match && err == nil {
		return
	}

	// Create the AST by parsing src.
	fset := token.NewFileSet() // positions are relative to fset
	node, err := parser.ParseFile(fset, file, nil, 0)
	if err != nil {
		fmt.Printf("parse file with error: %s\n", err)
		return
	}

	r = report{
		path: file,
	}

	ast.Inspect(node, func(n ast.Node) bool {
		if _, ok := n.(*ast.FuncDecl); ok {
			ast.Inspect(n, func(n ast.Node) bool {
				switch n := n.(type) {
				case *ast.AssignStmt, *ast.IncDecStmt:
					r.assignments++
				case *ast.CallExpr:
					r.branches++
				case *ast.IfStmt:
					if n.Else != nil {
						r.conditions++
					}
				case *ast.BinaryExpr, *ast.CaseClause:
					r.conditions++
				}
				return true
			})

			a := math.Pow(float64(r.assignments), 2)
			b := math.Pow(float64(r.branches), 2)
			c := math.Pow(float64(r.conditions), 2)

			r.score = int(math.Sqrt(a + b + c))
			return false
		}
		return true
	})

	return
}
