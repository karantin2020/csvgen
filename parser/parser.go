package parser

import (
	"go/ast"
	"go/parser"
	"go/token"
	"strings"
	"fmt"
	"os"
	"path"
)

type Parser struct {
	PkgPath     string
	PkgName     string
	StructNames []string
	AllStructs  bool
}

type visitor struct {
	*Parser

	name     string
	explicit bool
}

func (v *visitor) Visit(n ast.Node) (w ast.Visitor) {
	switch n := n.(type) {
	case *ast.Package:
		return v
	case *ast.File:
		v.PkgName = n.Name.String()
		return v

	case *ast.GenDecl:
		if !v.AllStructs {
			return nil
		}
		return v
	case *ast.TypeSpec:
		v.name = n.Name.String()

		// Allow to specify non-structs explicitly independent of '-all' flag.
		// if v.explicit {
		// 	v.StructNames = append(v.StructNames, v.name)
		// 	return nil
		// }
		return v
	case *ast.StructType:
		fmt.Printf("%s: %+v\n",v.name, n)
		// fmt.Println(n.Fields.List)
		for _, fl := range n.Fields.List {
			if e, ok := fl.Type.(*ast.StarExpr); !ok {
				fmt.Printf("%s %s\n", fl.Names[0], fl.Type)
			} else {
				fmt.Printf("%s *%s\n", fl.Names[0], e.X)
			}
		}
	
		v.StructNames = append(v.StructNames, v.name)
		return nil
	}
	return nil
}

func (p *Parser) Parse(fname string, isDir bool) error {
	var err error
	if p.PkgPath, err = getPkgPath(fname, isDir); err != nil {
		return err
	}

	fset := token.NewFileSet()
	if isDir {
		packages, err := parser.ParseDir(fset, fname, nil, parser.ParseComments)
		if err != nil {
			return err
		}

		for _, pckg := range packages {
			ast.Walk(&visitor{Parser: p}, pckg)
		}
	} else {
		f, err := parser.ParseFile(fset, fname, nil, parser.ParseComments)
		if err != nil {
			fmt.Println("Error parsing file:", err)
			return err
		}
		// fmt.Println(f)
		ast.Walk(&visitor{Parser: p}, f)
	}
	return nil
}

func getPkgPath(fname string, isDir bool) (string, error) {
	if !path.IsAbs(fname) {
		pwd, err := os.Getwd()
		if err != nil {
			return "", err
		}
		fname = path.Join(pwd, fname)
	}

	for _, p := range strings.Split(os.Getenv("GOPATH"), ":") {
		prefix := path.Join(p, "src") + "/"
		if rel := strings.TrimPrefix(fname, prefix); rel != fname {
			if !isDir {
				return path.Dir(rel), nil
			} else {
				return path.Clean(rel), nil
			}
		}
	}

	return "", fmt.Errorf("file '%v' is not in GOPATH", fname)
}

