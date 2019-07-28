package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"strings"

	astvisit "github.com/ungerik/go-astvisit"
	fs "github.com/ungerik/go-fs"
)

func main() {
	packageDir := fs.File("~/go/src/github.com/domonda/Domonda/pkg/db/postgresdb")

	packageDir.ListDir(func(file fs.File) error {
		if file.IsDir() || file.Ext() != ".go" || strings.HasSuffix(file.Name(), "_test.go") {
			return nil
		}
		src, err := file.ReadAll()
		if err != nil {
			return err
		}

		fset := token.NewFileSet()
		f, err := parser.ParseFile(fset, file.Name(), src, parser.ParseComments)
		if err != nil {
			return err
		}

		// ast.Print(fset, f)
		// os.Exit(0)

		visitor := &funcDeclVisitor{
			fset: fset,
			src:  src,
			offs: f.Pos(),
		}
		astvisit.Visit(f, visitor, nil)
		// err = ast.Print(fset, result)
		// if err != nil {
		// 	return err
		// }

		return nil
	})

}

type funcDeclVisitor struct {
	astvisit.VisitorImpl

	src  []byte
	fset *token.FileSet
	offs token.Pos
}

func (v *funcDeclVisitor) exprString(expr ast.Expr) string {
	return string(v.src[expr.Pos()-v.offs : expr.End()-v.offs])
}

func (v *funcDeclVisitor) VisitFuncDecl(funcDecl *ast.FuncDecl, cursor astvisit.Cursor) bool {
	if !funcDecl.Name.IsExported() {
		return false
	}

	resultErrorName, hasResultError := funcError(funcDecl)
	if !hasResultError {
		return false
	}

	if resultErrorName == "" {
		fmt.Println("Function error is not named:", v.fset.Position(funcDecl.Pos()))
		// fmt.Println(cursor.Path())
		fmt.Printf("%s(%s)\n", funcDecl.Name.Name, strings.Join(fieldNames(funcDecl.Type.Params.List), ", "))
		fmt.Println()
		return false
	}

	var (
		deferStmt   *ast.DeferStmt
		wrappedArgs []string
	)
	for _, stmt := range funcDecl.Body.List {
		deferStmt, wrappedArgs = asDeferWrapErrorStatement(stmt, resultErrorName)
		if deferStmt != nil {
			break
		}
	}

	params := fieldNames(funcDecl.Type.Params.List)

	if deferStmt == nil {
		fmt.Println("Function is missing error wrapping:", v.fset.Position(funcDecl.Pos()))
		// fmt.Println(cursor.Path())
		fmt.Printf("%s(%s)\n", funcDecl.Name.Name, strings.Join(params, ", "))
		fmt.Println()
		return false
	}

	if !equalStrings(params, wrappedArgs) {
		fmt.Println("Different functions args:", v.fset.Position(funcDecl.Pos()))
		fmt.Printf("%s(%s) <-> wrap(%s)\n", funcDecl.Name.Name, strings.Join(params, ", "), strings.Join(wrappedArgs, ", "))
		// fmt.Println(cursor.Path())
		// fmt.Printf("%s(%s)\n", funcDecl.Name.Name, strings.Join(params, ", "))
		fmt.Println()
		return false
	}

	// ast.Print(v.fset, stmt)
	// os.Exit(0)

	return false
}

func fieldNames(fields []*ast.Field) (names []string) {
	for _, field := range fields {
		for _, fieldName := range field.Names {
			names = append(names, fieldName.Name)
		}
	}
	return names
}

func equalStrings(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

func funcError(funcDecl *ast.FuncDecl) (name string, hasErr bool) {
	if funcDecl.Type.Results == nil {
		return "", false
	}
	results := funcDecl.Type.Results.List
	if len(results) == 0 {
		return "", false
	}
	last := results[len(results)-1]
	if id, ok := last.Type.(*ast.Ident); !ok || id.Name != "error" {
		return "", false
	}
	if len(last.Names) == 0 {
		return "", true
	}
	return last.Names[len(last.Names)-1].Name, true
}

func asDeferWrapErrorStatement(stmt ast.Stmt, resultErrorName string) (deferStmt *ast.DeferStmt, wrappedArgs []string) {
	deferStmt, ok := stmt.(*ast.DeferStmt)
	if !ok {
		return nil, nil
	}

	funSel, ok := deferStmt.Call.Fun.(*ast.SelectorExpr)
	if !ok {
		return nil, nil
	}

	if selectorExpString(funSel) != "wrap.ResultError" {
		return nil, nil
	}

	args := deferStmt.Call.Args

	if len(args) < 2 {
		return nil, nil
	}

	unary, ok := args[0].(*ast.UnaryExpr)
	if !ok || unary.Op != token.AND {
		return nil, nil
	}
	if ident, ok := unary.X.(*ast.Ident); !ok || ident.Name != resultErrorName {
		return nil, nil
	}

	_, ok = args[1].(*ast.BasicLit)
	if !ok {
		return nil, nil
	}

	for i := 2; i < len(args); i++ {
		argIdent, ok := args[i].(*ast.Ident)
		if !ok {
			wrappedArgs = append(wrappedArgs, fmt.Sprintf("XXX arg is not ast.Ident: %+v", args[i]))
			continue
		}
		wrappedArgs = append(wrappedArgs, argIdent.Name)
	}

	return deferStmt, wrappedArgs
}

func selectorExpString(sel *ast.SelectorExpr) string {
	return sel.X.(*ast.Ident).Name + "." + sel.Sel.Name
}

// t := field.Type

// ellipsis, varArgs := t.(*ast.Ellipsis)
// if varArgs {
// 	t = ellipsis.Elt
// }

// starExpr, isPtr := t.(*ast.StarExpr)
// if isPtr {
// 	t = starExpr.X
// }

// switch x := t.(type) {
// case *ast.Ident:
// case *ast.SelectorExpr:
// case *ast.ArrayType:
// case *ast.MapType:
// case *ast.FuncType:
// case *ast.InterfaceType:
// 	// fmt.Println("Interface methods:", x.Methods.NumFields())

// default:
// 	fmt.Printf("\n\nUNKNOWN: %#v %T\n\n\n", x, x)
// }
