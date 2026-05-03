package main

import (
	"bytes"
	"flag"
	"fmt"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"os"
	"strconv"
)

// go generate ./internal/services/... ./internal/repository/...
func main() {
	filePath := flag.String("file", "", "Go source file to rewrite")
	receiverType := flag.String("receiver", "", "receiver type to instrument")
	serviceName := flag.String("service", "", "service/interface name used in span names")
	flag.Parse()

	if *filePath == "" || *receiverType == "" || *serviceName == "" {
		fmt.Fprintln(os.Stderr, "usage: injecttrace -file <file.go> -receiver <receiverType> -service <ServiceName>")
		os.Exit(2)
	}

	if err := injectTracing(*filePath, *receiverType, *serviceName); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func injectTracing(filePath, receiverType, serviceName string) error {
	original, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("read file: %w", err)
	}

	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, filePath, original, parser.ParseComments)
	if err != nil {
		return fmt.Errorf("parse file: %w", err)
	}

	var instrumented bool

	for _, decl := range file.Decls {
		fn, ok := decl.(*ast.FuncDecl)
		if !ok || fn.Recv == nil || fn.Body == nil || len(fn.Recv.List) == 0 {
			continue
		}
		if exprName(fn.Recv.List[0].Type) != receiverType {
			continue
		}

		ctxName, ok := contextParamName(fn.Type.Params)
		if !ok || hasTracePrefix(fn, ctxName, serviceName) {
			continue
		}

		startCall, err := parser.ParseExpr(
			fmt.Sprintf(`otel.Tracer(%q).Start(%s, %q)`, serviceName, ctxName, serviceName+"."+fn.Name.Name),
		)
		if err != nil {
			return fmt.Errorf("build trace start for %s: %w", fn.Name.Name, err)
		}

		assign := &ast.AssignStmt{
			Lhs: []ast.Expr{ast.NewIdent(ctxName), ast.NewIdent("span")},
			Tok: token.DEFINE,
			Rhs: []ast.Expr{startCall},
		}
		deferEnd := &ast.DeferStmt{
			Call: &ast.CallExpr{
				Fun: &ast.SelectorExpr{
					X:   ast.NewIdent("span"),
					Sel: ast.NewIdent("End"),
				},
			},
		}

		fn.Body.List = append([]ast.Stmt{assign, deferEnd}, fn.Body.List...)
		instrumented = true
	}

	if !instrumented {
		return nil
	}

	ensureImport(file, "go.opentelemetry.io/otel")
	ast.SortImports(fset, file)

	var output bytes.Buffer
	if err := format.Node(&output, fset, file); err != nil {
		return fmt.Errorf("format file: %w", err)
	}

	if bytes.Equal(original, output.Bytes()) {
		return nil
	}

	if err := os.WriteFile(filePath, output.Bytes(), 0o644); err != nil {
		return fmt.Errorf("write file: %w", err)
	}

	return nil
}

func ensureImport(file *ast.File, importPath string) bool {
	for _, imp := range file.Imports {
		if importPathFor(imp) == importPath {
			return false
		}
	}

	newSpec := &ast.ImportSpec{
		Path: &ast.BasicLit{Kind: token.STRING, Value: strconv.Quote(importPath)},
	}

	for _, decl := range file.Decls {
		gen, ok := decl.(*ast.GenDecl)
		if ok && gen.Tok == token.IMPORT {
			gen.Specs = append(gen.Specs, newSpec)
			file.Imports = append(file.Imports, newSpec)
			return true
		}
	}

	file.Decls = append([]ast.Decl{
		&ast.GenDecl{
			Tok:   token.IMPORT,
			Specs: []ast.Spec{newSpec},
		},
	}, file.Decls...)
	file.Imports = append(file.Imports, newSpec)
	return true
}

func contextParamName(params *ast.FieldList) (string, bool) {
	if params == nil {
		return "", false
	}

	for _, field := range params.List {
		if !isContextType(field.Type) || len(field.Names) == 0 {
			continue
		}
		if field.Names[0].Name == "_" {
			field.Names[0].Name = "ctx"
		}
		return field.Names[0].Name, true
	}

	return "", false
}

func hasTracePrefix(fn *ast.FuncDecl, ctxName, serviceName string) bool {
	if len(fn.Body.List) < 2 {
		return false
	}

	assign, ok := fn.Body.List[0].(*ast.AssignStmt)
	if !ok || assign.Tok != token.DEFINE || len(assign.Lhs) != 2 || len(assign.Rhs) != 1 {
		return false
	}

	if identName(assign.Lhs[0]) != ctxName || identName(assign.Lhs[1]) != "span" {
		return false
	}

	call, ok := assign.Rhs[0].(*ast.CallExpr)
	if !ok {
		return false
	}

	startSelector, ok := call.Fun.(*ast.SelectorExpr)
	if !ok || startSelector.Sel.Name != "Start" {
		return false
	}

	if len(call.Args) < 2 || stringLiteralValue(call.Args[1]) != serviceName+"."+fn.Name.Name {
		return false
	}

	deferStmt, ok := fn.Body.List[1].(*ast.DeferStmt)
	if !ok {
		return false
	}

	endSelector, ok := deferStmt.Call.Fun.(*ast.SelectorExpr)
	return ok && identName(endSelector.X) == "span" && endSelector.Sel.Name == "End"
}

func isContextType(expr ast.Expr) bool {
	selector, ok := expr.(*ast.SelectorExpr)
	return ok && identName(selector.X) == "context" && selector.Sel.Name == "Context"
}

func exprName(expr ast.Expr) string {
	switch typed := expr.(type) {
	case *ast.Ident:
		return typed.Name
	case *ast.StarExpr:
		return exprName(typed.X)
	default:
		return ""
	}
}

func identName(expr ast.Expr) string {
	ident, ok := expr.(*ast.Ident)
	if !ok {
		return ""
	}
	return ident.Name
}

func importPathFor(spec *ast.ImportSpec) string {
	if spec == nil || spec.Path == nil {
		return ""
	}
	value, err := strconv.Unquote(spec.Path.Value)
	if err != nil {
		return spec.Path.Value
	}
	return value
}

func stringLiteralValue(expr ast.Expr) string {
	lit, ok := expr.(*ast.BasicLit)
	if !ok || lit.Kind != token.STRING {
		return ""
	}
	value, err := strconv.Unquote(lit.Value)
	if err != nil {
		return lit.Value
	}
	return value
}
