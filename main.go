package main

import (
	"flag"
	"fmt"
	"go/parser"
	"go/token"
	"log"
	"os"
	"strings"
)

const (
	binaryName    = "go-mock"
	binaryVersion = "v0.1.1"
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
	filename := os.Getenv("GOFILE")
	g.parseFile(filename, targetInterfaces)
	g.generateFiles()
}

type Generator struct {
	files []File
}

func (g *Generator) parseFile(filename string, targets map[string]string) {
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, filename, nil, 0)
	if err != nil {
		log.Fatal(err)
	}
	files := parseFiles(f, filename, targets)
	g.files = append(g.files, files...)
}

func (g *Generator) generateFiles() {
	wd, _ := os.Getwd()
	for _, f := range g.files {
		line := os.Getenv("GOLINE")
		log.Printf("generating file: %s/%s:L%s\n", wd, f.filename, line)
		f.generate()
	}
}
