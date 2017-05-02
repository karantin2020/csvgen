package main

import (
	"os"
	"github.com/karantin2020/csvgen/parser"
)

func main() {
	fname := "./fixture/test.go"
	fInfo, err := os.Stat(fname)
	if err != nil {
		return
	}
	p := parser.Parser{AllStructs: true}
	if err := p.Parse(fname, fInfo.IsDir()); err != nil {
		return
	}
}