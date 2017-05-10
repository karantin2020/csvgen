package main

import (
	"fmt"
	"github.com/karantin2020/csvgen/parser"
	"os"
)

var (
	p = parser.Parser{AllStructs: true}
)

func main() {
	fname := "./fixture"
	fInfo, err := os.Stat(fname)
	if err != nil {
		return
	}
	// p := parser.Parser{AllStructs: true}
	if err := p.Parse(fname, fInfo.IsDir()); err != nil {
		return
	}
	fmt.Println("Parse Error:", p.Error)
	fmt.Println(p.StructMap)
	fmt.Printf("%+v\n", p.Structs)
}
