package main

import (
	"github.com/reactive-futakotamagawa/isu-isu-h/tools/intro-h/pprotein/echo"
	"github.com/reactive-futakotamagawa/isu-isu-h/tools/intro-h/pprotein/integrate"
	"golang.org/x/tools/go/analysis/multichecker"
)

func main() { multichecker.Main(integrate.IntegrateMain, echo.Analyzer) }
