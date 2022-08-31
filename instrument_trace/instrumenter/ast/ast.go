package ast

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"

	"golang.org/x/tools/go/ast/astutil"
)

type instrumenter struct {
	traceImport string
	tracePkg    string
	traceFunc   string
}

func New(traceImport, tracePkg, traceFunc string) *instrumenter {
	return &instrumenter{
		traceImport: traceImport,
		tracePkg:    tracePkg,
		traceFunc:   traceFunc,
	}
}

func hasFuncDecl(f *ast.File) bool {
	if len(f.Decls) == 0 {
		return false
	}

	for _, decl := range f.Decls {
		_, ok := decl.(*ast.FuncDecl)
		if ok {
			return true
		}
	}

	return false
}

func (a instrumenter) Instrument(filename string) ([]byte, error) {
	fset := token.NewFileSet()
	curAST, err := parser.ParseFile(fset, filename, nil, parser.ParseComments)
	if err != nil {
		return nil, fmt.Errorf("error parsing %s: %w", parser.ParseComments)
	}

	if !hasFuncDecl(curAST) {
		return nil, nil
	}

	// 在 AST 上添加包导入语句
	astutil.AddImport(fset, curAST, a.traceImport)

	// 向 AST 上的所有函数注入 Trace 函数
	a.addDeferTraceIntoFuncDecls(curAST)

	buf := &bytes.Buffer{}
	err = format.Node(buf, fset, curAST) // 将修改后的 AST 转换回 Go 源码
	if err != nil {
		return nil, fmt.Errorf("error formatting new code: %w", err)
	}
	return buf.Bytes(), nil // 返回转换后的 Go 源码
}

func (a instrumenter) addDeferTraceIntoFuncDecls(f *ast.File) {
	for _, decl := range f.Decls { // 遍历所有声明语句
		fd, ok := decl.(*ast.FuncDecl) // 类型断言： 是否为函数声明
		if ok {
			// 如果是函数声明，则注入跟踪设施
			a.addDeferStmt(fd)
		}
	}
}

func (a instrumenter) addDeferStmt(fd *ast.FuncDecl) (added bool) {
	stmts := fd.Body.List

	// 判断 "defer trace.Trace()()" 语句是否存在
	for _, stmt := range stmts {
		ds, ok := stmt.(*ast.DeferStmt)
		if !ok {
			// 如果不是 defer 语句，则继续 for 循环
			continue
		}

		// 如果是 defer 语句，则要进一步判断是否为 defer trace.Trace()()
		ce, ok := ds.Call.Fun.(*ast.CallExpr)
		if !ok {
			continue
		}

		se, ok := ce.Fun.(*ast.SelectorExpr)
		if !ok {
			continue
		}

		x, ok := se.X.(*ast.Ident)
		if !ok {
			continue
		}
		if (x.Name == a.tracePkg) && (se.Sel.Name == a.traceFunc) {
			// defer trace.Trace()() 存在，返回
			return false
		}
	}

	// 没有找到 "defer trace.Trace()()". 注入一个新的跟踪语句
	// 在 AST 上构造一个 defer trace.Trace()()
	ds := &ast.DeferStmt{
		Call: &ast.CallExpr{
			Fun: &ast.CallExpr{
				Fun: &ast.SelectorExpr{
					X: &ast.Ident{
						Name: a.tracePkg,
					},
					Sel: &ast.Ident{
						Name: a.traceFunc,
					},
				},
			},
		},
	}

	newList := make([]ast.Stmt, len(stmts)+1)
	copy(newList[1:], stmts)
	newList[0] = ds // 注入新构造的 defer 语句
	fd.Body.List = newList
	return true
}
