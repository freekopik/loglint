package main

import (
	"github.com/freekopik/loglint/pkg/loglint"
	"golang.org/x/tools/go/analysis"
)

type analyzerPlugin struct{}

func (*analyzerPlugin) GetAnalyzers() []*analysis.Analyzer {
	return []*analysis.Analyzer{
		loglint.Analyzer,
	}
}

var AnalyzerPlugin analyzerPlugin
