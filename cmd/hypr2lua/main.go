package main

import (
	"fmt"
	"log"
	"os"

	"github.com/phillezi/hypr2lua/pkg/hypr/ast"
	"github.com/phillezi/hypr2lua/pkg/hypr/compiler"
	"github.com/phillezi/hypr2lua/pkg/lua/stubs"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: hypr2lua <config>")
		os.Exit(1)
	}

	inputPath := os.Args[1]

	f, err := os.Open(inputPath)
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
	defer f.Close()

	parser := ast.NewParser(f)
	astfile, err := parser.Parse()
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	schema, err := stubs.NewLoader().LoadDir("/usr/share/hypr/stubs")
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	cc := compiler.New(schema)
	if err := cc.Compile(astfile); err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
}
