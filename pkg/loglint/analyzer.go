package loglint

import (
	"go/ast"
	"path/filepath"
	"regexp"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

var Analyzer = &analysis.Analyzer{
	Name:     "loglint",
	Doc:      "checks log messages for style and security",
	Run:      run,
	Requires: []*analysis.Analyzer{inspect.Analyzer},
}

func run(pass *analysis.Pass) (interface{}, error) {
	inspect := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)
	cfg := getConfig()

	var excludeRes []*regexp.Regexp
	for _, pattern := range cfg.Settings.ExcludeFiles {
		re, err := regexp.Compile("(?i)" + pattern)
		if err == nil {
			excludeRes = append(excludeRes, re)
		}
	}

	nodeFilter := []ast.Node{(*ast.CallExpr)(nil)}

	inspect.Preorder(nodeFilter, func(n ast.Node) {
		pos := pass.Fset.Position(n.Pos())
		if !pos.IsValid() {
			return
		}

		absPath, _ := filepath.Abs(pos.Filename)
		fullPath := filepath.ToSlash(absPath)
		baseName := filepath.Base(fullPath)

		for _, re := range excludeRes {
			if re.MatchString(fullPath) || re.MatchString(baseName) {
				return
			}
		}

		call := n.(*ast.CallExpr)
		if !isLogCall(call, cfg.TargetMethods) {
			return
		}

		if len(call.Args) > 0 {
			checkLogArgument(pass, call.Args[0], cfg)
		}
	})

	return nil, nil
}

func isLogCall(call *ast.CallExpr, targetMethods []string) bool {
	sel, ok := call.Fun.(*ast.SelectorExpr)
	if !ok {
		return false
	}

	methodName := sel.Sel.Name
	for _, m := range targetMethods {
		if m == methodName {
			if ident, ok := sel.X.(*ast.Ident); ok {
				pkg := ident.Name
				if pkg == "log" || pkg == "slog" || pkg == "zap" || pkg == "logger" {
					return true
				}
			}
			break
		}
	}
	return false
}
