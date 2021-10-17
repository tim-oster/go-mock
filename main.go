package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"golang.org/x/tools/go/packages"
)

const (
	binaryName    = "go-mock"
	binaryVersion = "v0.1.0"
)

var (
	keepCtx    = flag.Bool("keepctx", false, "does not remove ctx parameters when present as first param")
	unexported = flag.Bool("unexported", false, "generates an unexported mock struct")
)

func Usage() {
	fmt.Fprintf(os.Stderr, "Usage of %s:\n", binaryName)
	fmt.Fprintf(os.Stderr, "\t%s [flags] Pattern[:rename][,...]\n", binaryName)
	fmt.Fprintf(os.Stderr, "Flags:\n")
	flag.PrintDefaults()
}

func main() {
	log.SetFlags(0)
	log.SetPrefix("go-mock: ")
	flag.Usage = Usage
	flag.Parse()

	patternArg := flag.Arg(0)
	if len(patternArg) == 0 {
		fmt.Fprintf(os.Stderr, "missing pattern argument\n\n")
		flag.Usage()
		os.Exit(1)
	}

	targetInterfaces := map[string]string{}
	patterns := strings.Split(patternArg, ",")
	for _, p := range patterns {
		parts := strings.SplitN(p, ":", 2)
		var rename string
		if len(parts) == 2 {
			rename = parts[1]
		}
		targetInterfaces[parts[0]] = rename
	}

	var g Generator
	g.parsePackage(targetInterfaces)
	g.generateFiles()
}

type Generator struct {
	files []File
}

func (g *Generator) parsePackage(targets map[string]string) {
	cfg := &packages.Config{Mode: packages.NeedSyntax | packages.NeedTypes | packages.NeedFiles | packages.NeedImports}
	pkgs, err := packages.Load(cfg, ".")
	if err != nil {
		log.Fatal(err)
	}
	packages.PrintErrors(pkgs)
	for _, pkg := range pkgs {
		files := parseFiles(pkg, targets)
		g.files = append(g.files, files...)
	}
}

func (g *Generator) generateFiles() {
	for _, f := range g.files {
		log.Printf("generating file: %s\n", f.filename)
		f.generate()
	}
}
