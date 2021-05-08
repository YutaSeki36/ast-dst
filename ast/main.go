package main

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/parser"
	"go/printer"
	"go/token"
	"os"
	"strconv"
)

func main() {
	src := `
package main

import "fmt"

// メソッド1
func call1() {
	fmt.Print("1")
}

// メソッド2
func call2() {
	fmt.Print("2")
}

// main メイン処理
func main() {
  fmt.Println("Hello, World!")
}
`
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, "", src, 4)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	ast.Inspect(f, func(n ast.Node) bool {

		switch nc := n.(type) {
		case *ast.FuncDecl:
			if nc.Name.Name == "call1" {
				nc.Pos()
				nc.Body = &ast.BlockStmt{
					List: []ast.Stmt{
						&ast.ExprStmt{
							X: &ast.CallExpr{
								Fun: &ast.SelectorExpr{
									X:   ast.NewIdent("fmt"),
									Sel: ast.NewIdent("Println"),
								},
								Args: []ast.Expr{
									&ast.BasicLit{
										Kind:  token.STRING,
										Value: strconv.Quote("一行目追加"),
									},
								},
							},
						},
						&ast.ExprStmt{
							X: &ast.CallExpr{
								Fun: &ast.SelectorExpr{
									X:   ast.NewIdent("fmt"),
									Sel: ast.NewIdent("Println"),
								},
								Args: []ast.Expr{
									&ast.BasicLit{
										Kind:  token.STRING,
										Value: strconv.Quote("二行目追加"),
									},
								},
							},
						},
					},
				}
			}
		}

		return true
	})

	ast.Print(fset, f)
	var output []byte
	buf := bytes.NewBuffer(output)
	if err := printer.Fprint(buf, fset, f); err != nil {
		fmt.Errorf(err.Error())
		return
	}

	fmt.Printf("%s", buf.String())
}
