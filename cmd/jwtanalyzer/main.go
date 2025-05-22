package main

import (
	"github.com/muizu555/jwt-analyzer/analyzer"
	"golang.org/x/tools/go/analysis/singlechecker"
)

func main() {
	singlechecker.Main(analyzer.IssuerAnalyzer)
}
