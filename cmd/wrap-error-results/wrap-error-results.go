package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/printer"
	"go/token"
	"os/exec"
	"strings"

	astvisit "github.com/ungerik/go-astvisit"
	fs "github.com/ungerik/go-fs"
)

func main() {
	packageDir := fs.File("~/go/src/github.com/domonda/Domonda/pkg/db/postgresdb")

	fset := token.NewFileSet()

	packageDir.ListDir(func(file fs.File) error {
		if file.IsDir() || file.Ext() != ".go" || strings.HasSuffix(file.Name(), "_test.go") {
			return nil
		}
		src, err := file.ReadAll()
		if err != nil {
			return err
		}

		f, err := parser.ParseFile(fset, file.Name(), src, parser.ParseComments)
		if err != nil {
			return err
		}

		// ast.Print(fset, f)
		// os.Exit(0)

		hasChanges := false

		visitor := &funcDeclVisitor{
			fset: fset,

			onFuncCorrectlyWrapped: func(funcDecl *ast.FuncDecl, cursor astvisit.Cursor, deferStmt *ast.DeferStmt) {
				// ast.Print(fset, deferStmt)
			},

			onFuncNotWrapped: func(problem string, funcDecl *ast.FuncDecl, cursor astvisit.Cursor) {
				fmt.Printf("%s\n\t%s\n", problem, fset.Position(funcDecl.Pos()))
				fmt.Printf("\t%s(%s)\n", funcDecl.Name.Name, strings.Join(fieldNames(funcDecl.Type.Params.List), ", "))
			},

			onFuncWronglyWrapped: func(problem string, wrappedArgs []string, funcDecl *ast.FuncDecl, cursor astvisit.Cursor, deferStmt *ast.DeferStmt, deferStmtIndex int) {
				args := fieldNames(funcDecl.Type.Params.List)
				fmt.Printf("%s\n\t%s\n", problem, fset.Position(funcDecl.Pos()))
				fmt.Printf("\t%s(%s) <-> wrap(%s)\n", funcDecl.Name.Name, strings.Join(args, ", "), strings.Join(wrappedArgs, ", "))

				funcDecl.Body.List[deferStmtIndex] = newWrapDeferStmt(funcDecl)
				hasChanges = true
			},
		}

		result := astvisit.Visit(f, visitor, nil)

		if hasChanges {
			fmt.Println("Rewriting file:", file.LocalPath())

			newFile := file // file.Dir().Join("out_" + file.Name())
			writer, err := newFile.OpenWriter()
			if err != nil {
				return err
			}
			defer writer.Close()

			err = printer.Fprint(writer, fset, result)
			if err != nil {
				return err
			}

			out, err := exec.Command("gofmt", "-w", newFile.LocalPath()).CombinedOutput()
			if err != nil {
				fmt.Println("gofmt:", string(out))
				return err
			}
		}

		return nil
	})
}

type funcDeclVisitor struct {
	astvisit.VisitorImpl

	fset                   *token.FileSet
	onFuncCorrectlyWrapped func(funcDecl *ast.FuncDecl, cursor astvisit.Cursor, deferStmt *ast.DeferStmt)
	onFuncNotWrapped       func(problem string, funcDecl *ast.FuncDecl, cursor astvisit.Cursor)
	onFuncWronglyWrapped   func(problem string, wrappedArgs []string, funcDecl *ast.FuncDecl, cursor astvisit.Cursor, deferStmt *ast.DeferStmt, deferStmtIndex int)
}

func (v *funcDeclVisitor) VisitFuncDecl(funcDecl *ast.FuncDecl, cursor astvisit.Cursor) bool {
	if !funcDecl.Name.IsExported() {
		return false
	}

	// TODO find wrap.ResultError without defer

	resultErrorName, hasResultError := funcError(funcDecl)
	if !hasResultError {
		return false
	}

	if resultErrorName == "" {
		v.onFuncNotWrapped("Function error result is not named", funcDecl, cursor)
		return false
	}

	var (
		deferStmt      *ast.DeferStmt
		deferStmtIndex int
		wrappedArgs    []string
	)
	for i, stmt := range funcDecl.Body.List {
		deferStmt, wrappedArgs = asDeferWrapErrorStatement(stmt, resultErrorName)
		if deferStmt != nil {
			deferStmtIndex = i
			break
		}
	}

	if deferStmt == nil {
		v.onFuncNotWrapped("Function is missing error wrapping", funcDecl, cursor)
		return false
	}

	args := fieldNames(funcDecl.Type.Params.List)
	if !equalStrings(args, wrappedArgs) {
		v.onFuncWronglyWrapped("Different functions args", wrappedArgs, funcDecl, cursor, deferStmt, deferStmtIndex)
		return false
	}

	// ast.Print(v.fset, stmt)
	// os.Exit(0)

	v.onFuncCorrectlyWrapped(funcDecl, cursor, deferStmt)
	return false
}

func newWrapDeferStmt(funcDecl *ast.FuncDecl) *ast.DeferStmt {
	paramNames := fieldNames(funcDecl.Type.Params.List)
	errName, hasErr := funcError(funcDecl)
	if errName == "" || !hasErr {
		panic("function has not named error result")
	}
	s := &ast.DeferStmt{
		Call: &ast.CallExpr{
			Fun: &ast.SelectorExpr{X: &ast.Ident{Name: "wrap"}, Sel: &ast.Ident{Name: "ResultError"}},
			Args: []ast.Expr{
				&ast.UnaryExpr{Op: token.AND, X: &ast.Ident{Name: errName}},
				&ast.BasicLit{Kind: token.STRING, Value: `"` + funcDecl.Name.Name + `"`},
			},
		},
	}
	for _, arg := range paramNames {
		s.Call.Args = append(s.Call.Args, &ast.Ident{Name: arg})
	}
	return s
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
