package loglint

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/printer"
	"go/token"
	"strconv"
	"strings"
	"unicode"

	"golang.org/x/tools/go/analysis"
)

func checkLogArgument(pass *analysis.Pass, arg ast.Expr, cfg Config) {
	argCode := renderNode(pass.Fset, arg)
	argCodeLower := strings.ToLower(argCode)

	for key, replacement := range cfg.SensitiveWords {
		if strings.Contains(argCodeLower, strings.ToLower(key)) {
			pass.Report(analysis.Diagnostic{
				Pos:     arg.Pos(),
				End:     arg.End(),
				Message: "log message contains potentially sensitive data",
				SuggestedFixes: []analysis.SuggestedFix{
					{
						Message: fmt.Sprintf("replace with safe message: %s", replacement),
						TextEdits: []analysis.TextEdit{{
							Pos:     arg.Pos(),
							End:     arg.End(),
							NewText: []byte(strconv.Quote(replacement)),
						}},
					},
				},
			})
			return
		}
	}

	if lit, ok := arg.(*ast.BasicLit); ok && lit.Kind == token.STRING {
		text, _ := strconv.Unquote(lit.Value)

		if !startsWithLowercase(text) && len(text) > 0 {
			runes := []rune(text)
			runes[0] = unicode.ToLower(runes[0])
			pass.Report(analysis.Diagnostic{
				Pos:     arg.Pos(),
				End:     arg.End(),
				Message: "log message must start with a lowercase letter",
				SuggestedFixes: []analysis.SuggestedFix{
					{
						Message: "convert to lowercase",
						TextEdits: []analysis.TextEdit{{
							Pos:     arg.Pos(),
							End:     arg.End(),
							NewText: []byte(strconv.Quote(string(runes))),
						}},
					},
				},
			})
		}

		if !isEnglish(text) {
			pass.Reportf(arg.Pos(), "log message must be in English only")
		}

		if hasSpecChars(text, cfg.Settings.AllowedSymbols) {
			fixed := cleanSpecChars(text, cfg.Settings.AllowedSymbols)
			pass.Report(analysis.Diagnostic{
				Pos:     arg.Pos(),
				End:     arg.End(),
				Message: "log message must not contain special characters",
				SuggestedFixes: []analysis.SuggestedFix{
					{
						Message: "remove forbidden characters",
						TextEdits: []analysis.TextEdit{{
							Pos:     arg.Pos(),
							End:     arg.End(),
							NewText: []byte(strconv.Quote(fixed)),
						}},
					},
				},
			})
		}
	}
}

func hasSpecChars(t string, allowed string) bool {
	for _, r := range t {
		if !unicode.IsLetter(r) && !unicode.IsNumber(r) && !unicode.IsSpace(r) && !strings.ContainsRune(allowed, r) {
			return true
		}
	}
	return false
}

func cleanSpecChars(t string, allowed string) string {
	var sb strings.Builder
	for _, r := range t {
		if unicode.IsLetter(r) || unicode.IsNumber(r) || unicode.IsSpace(r) || strings.ContainsRune(allowed, r) {
			sb.WriteRune(r)
		}
	}
	return strings.TrimSpace(sb.String())
}

func startsWithLowercase(t string) bool {
	t = strings.TrimSpace(t)
	if t == "" {
		return true
	}
	runes := []rune(t)
	if !unicode.IsLetter(runes[0]) {
		return true
	}
	return unicode.IsLower(runes[0])
}

func isEnglish(t string) bool {
	for _, r := range t {
		if r > unicode.MaxASCII {
			return false
		}
	}
	return true
}

func renderNode(fset *token.FileSet, node ast.Node) string {
	var buf bytes.Buffer
	if err := printer.Fprint(&buf, fset, node); err != nil {
		return ""
	}
	return buf.String()
}
