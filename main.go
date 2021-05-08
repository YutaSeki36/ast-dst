package main

import (
	"bytes"
	"fmt"
	"go/printer"
	"go/token"
	"os"
	"strconv"

	"github.com/dave/dst"
	"github.com/dave/dst/decorator"
)

func main() {
	src := `
package main

import "fmt"

// 関数1
func call1() {
	fmt.Print("1")
}

// 関数2
func call2() {
	fmt.Print("2")
}

// main メイン処理
func main() {
  fmt.Println("Hello, World!")
}
`
	fset := token.NewFileSet()
	f, err := decorator.ParseFile(fset, "", src, 4)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	dst.Inspect(f, func(n dst.Node) bool {

		switch nc := n.(type) {
		case *dst.FuncDecl:
			if nc.Name.Name == "call1" {
				fmt.Println(nc.Recv)
				nc.Body = &dst.BlockStmt{
					List: []dst.Stmt{
						&dst.ExprStmt{
							X: &dst.CallExpr{
								Fun: &dst.SelectorExpr{
									X:   dst.NewIdent("fmt"),
									Sel: dst.NewIdent("Println"),
								},
								Args: []dst.Expr{
									&dst.BasicLit{
										Kind:  token.STRING,
										Value: strconv.Quote("一行目追加"),
									},
								},
							},
						},
						&dst.ExprStmt{
							X: &dst.CallExpr{
								Fun: &dst.SelectorExpr{
									X:   dst.NewIdent("fmt"),
									Sel: dst.NewIdent("Println"),
								},
								Args: []dst.Expr{
									&dst.BasicLit{
										Kind:  token.STRING,
										Value: strconv.Quote("2行目追加"),
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

	decorator.Print(f)
	var output []byte
	buf := bytes.NewBuffer(output)
	if err := printer.Fprint(buf, fset, f); err != nil {
		fmt.Errorf(err.Error())
		return
	}
	fmt.Printf("%s", buf.String())
}
