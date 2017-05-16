package main

import (
	"fmt"
	"github.com/karantin2020/csvgen/parser"
	"os"
)

var (
	pp = parser.Parser{AllStructs: true}
)

func main() {
	fname := "./fixture"
	fInfo, err := os.Stat(fname)
	if err != nil {
		return
	}
	// pp := parser.Parser{AllStructs: true}
	if err := pp.Parse(fname, fInfo.IsDir()); err != nil {
		return
	}
	fmt.Println("Parse Error:", pp.Error)
	fmt.Println(pp.StructMap)
	fmt.Printf("%+v\n", pp.Structs)
}
