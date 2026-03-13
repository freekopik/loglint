package main

import (
	"github.com/freekopik/loglint/pkg/loglint"
	"golang.org/x/tools/go/analysis/singlechecker"
)

func main() {
	singlechecker.Main(loglint.Analyzer)
}
