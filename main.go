package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"math"
	"os"
	"path/filepath"

	"github.com/jamieabc/go-abc/order"

	"github.com/jamieabc/go-abc/report"
)

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

	var reports []report.Report
	for _, f := range files {
		reports = append(reports, analyze(f))
	}
	byScore := order.NewOrderByScore(reports)
	byScore.Sort()
	fmt.Printf("%s\n", byScore.Report())
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
				if isIgnored(path) {
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
		if isIgnored(root) {
			return
		}
		files = append(files, root)
	}

	return
}

func isIgnored(path string) bool {
	goMatch, err := filepath.Match("*.go", filepath.Base(path))
	testMatch, err := filepath.Match("*_test.go", filepath.Base(path))
	if goMatch && !testMatch && err == nil {
		return false
	}
	return true
}

func analyze(file string) (r report.Report) {
	// Create the AST by parsing src.
	fset := token.NewFileSet() // positions are relative to fset
	node, err := parser.ParseFile(fset, file, nil, 0)
	if err != nil {
		fmt.Printf("parse file with error: %s\n", err)
		return
	}

	r = report.Report{
		Path: file,
	}

	ast.Inspect(node, func(n ast.Node) bool {
		if _, ok := n.(*ast.FuncDecl); ok {
			ast.Inspect(n, func(n ast.Node) bool {
				switch n := n.(type) {
				case *ast.AssignStmt, *ast.IncDecStmt:
					r.ABC.Assignments++
				case *ast.CallExpr:
					r.ABC.Branches++
				case *ast.IfStmt:
					if n.Else != nil {
						r.ABC.Conditions++
					}
				case *ast.BinaryExpr, *ast.CaseClause:
					r.ABC.Conditions++
				}
				return true
			})

			a := math.Pow(float64(r.ABC.Assignments), 2)
			b := math.Pow(float64(r.ABC.Branches), 2)
			c := math.Pow(float64(r.ABC.Conditions), 2)

			r.Score = int(math.Sqrt(a + b + c))
			return false
		}
		return true
	})

	return
}
