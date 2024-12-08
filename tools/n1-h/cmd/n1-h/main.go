package main

import (
	n1h "github.com/reactive-futakotamagawa/isu-isu-h/tools/n1-h"
	"golang.org/x/tools/go/analysis/unitchecker"
)

func main() { unitchecker.Main(n1h.Analyzer) }
