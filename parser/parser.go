package parser

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path"
	"strings"
)

type FieldInfo struct {
	Name string
	Type string
}

type StructInfo struct {
	Name   string
	Fields []FieldInfo
}

type Parser struct {
	PkgPath     string
	PkgName     string
	StructNames []string
	AllStructs  bool

	StructMap map[string]map[string]string
	Structs   []StructInfo
	Error     bool
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

		return v
	case *ast.StructType:
		var tempStr map[string]string

		// Check if that struct name is cached already
		if _, prs := v.Parser.StructMap[v.name]; prs {
			// If cached then set Error true
			v.Parser.Error = true
		} else {
			// If not cached then parse struct fields
			v.Parser.StructMap[v.name] = make(map[string]string)
			tempStr = v.Parser.StructMap[v.name]

			parsedStruct := StructInfo{v.name, []FieldInfo{}}

			for _, fl := range n.Fields.List {
				tmpField := FieldInfo{}
				vname := fmt.Sprintf("%s", fl.Names[0])
				prsf := false
				// Check if struct field name exists
				if _, prsf = tempStr[vname]; prsf {
					// If exists then set Error true
					v.Parser.Error = true
				} else {
					// If not then set field name
					tmpField.Name = vname

					// Check for pointer type
					if e, ok := fl.Type.(*ast.StarExpr); ok {
						tmpField.Type = fmt.Sprintf("*%s", e.X)
					} else {
						tmpField.Type = fmt.Sprintf("%s", fl.Type)

					}
					// Set field type
					tempStr[tmpField.Name] = tmpField.Type
					parsedStruct.Fields = append(parsedStruct.Fields, tmpField)
				}
			}

			v.Structs = append(v.Structs, parsedStruct)
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
	p.StructMap = make(map[string]map[string]string)
	p.Error = false

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
		ast.Walk(&visitor{Parser: p}, f)
	}
	// fmt.Println(p.PkgName)
	// fmt.Println(p.PkgPath)

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
