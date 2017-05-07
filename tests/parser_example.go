package main

import (
	"os"
	"fmt"
	"github.com/karantin2020/csvgen/parser"
)

var (
	p = parser.Parser{AllStructs: true}
)

func main() {
	fname := "./fixture/test.go"
	fInfo, err := os.Stat(fname)
	if err != nil {
		return
	}
	// p := parser.Parser{AllStructs: true}
	if err := p.Parse(fname, fInfo.IsDir()); err != nil {
		return
	}
	fmt.Println(p.Error)
	fmt.Println(p.StructMap)
}