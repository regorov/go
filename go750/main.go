package main

import (
	"fmt"
	"github.com/kierdavis/ansi"
	"github.com/kierdavis/argparse"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
)

var Output = os.Stdout

func emit(name string, operands ...Operand) {
	fmt.Fprintln(Output, NewInstruction(name, operands...).String())
}

func emitLabel(name string) {
	fmt.Fprintln(Output, name+":")
}

func emitHeader() {
	emit("mov", PC, Label("main"))
}

func emitFooter() {

}

type Args struct {
	Dir string
}

func main() {
	p := argparse.New("A minimal Go compiler for K750")
	p.Argument("Dir", 1, argparse.Store, "dir", "The package directory to compile. It should contain a 'main' package.")

	args := &Args{}
	err := p.Parse(args)
	if err != nil {
		if cmdLineErr, ok := err.(argparse.CommandLineError); ok {
			ansi.Fprintln(os.Stderr, ansi.RedBold, string(cmdLineErr))
			p.Usage()
			os.Exit(2)

		} else {
			ansi.Fprintf(os.Stderr, ansi.RedBold, "Error: %s\n", err.Error())
			os.Exit(1)
		}
	}

	fset := token.NewFileSet()
	pkgs, err := parser.ParseDir(fset, args.Dir, nil, parser.DeclarationErrors)
	if err != nil {
		ansi.Fprintf(os.Stderr, ansi.RedBold, "Error: %s\n", err.Error())
		os.Exit(1)
	}

	pkg, ok := pkgs["main"]
	if !ok {
		ansi.Fprintf(os.Stderr, ansi.RedBold, "Error: main package was not found.")
		os.Exit(1)
	}

	file := ast.MergePackageFiles(pkg, ast.FilterFuncDuplicates|ast.FilterImportDuplicates)

	emitHeader()
	defer emitFooter()

	ast.Walk(NewFileVisitor(), file)
}
