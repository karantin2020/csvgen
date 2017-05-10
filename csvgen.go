package main

import (
	"fmt"
	. "github.com/dave/jennifer/jen"
	"github.com/karantin2020/cli"
	"github.com/karantin2020/csvgen/parser"
	"log"
	"os"
	"path/filepath"
	"strings"
)

var (
	// Config vars
	pkg    string
	subpkg string
	out    string
	fname  string
	fInfo  os.FileInfo
	// apn bool

	// Package vars
	pkgCnt string
	f      *File
	p      = parser.Parser{AllStructs: true}
)

func main() {
	flags := cli.New("This app generates csv Marshall and Unmarshal functions", "0.0.1")
	flags.StringVarP(&pkg, "pkg", "p", "", "output package")
	flags.StringVarP(&subpkg, "subpkg", "s", "", "output subpkg name")
	flags.StringVarP(&fname, "fname", "f", "", "input file")
	flags.StringVarP(&out, "out", "o", "", "output file")
	// flags.BoolVarP(&apn, "append", "a", false, "append result of template render to file or not")

	flags.Parse()

	pkgCnt = "main"
	if pkg != "" {
		pkgCnt = pkg
	}
	if subpkg != "" {
		pkgCnt = subpkg
		pkg = ""
	}

	if fname == "" {
		log.Fatal("didn't pass source file to parse")
		os.Exit(1)
	}
	var err error
	if fInfo, err = os.Stat(fname); err != nil {
		fmt.Println("Couldn't find source file to parse.", err)
		os.Exit(1)
	}

	// if out == "" {
	// 	// log.Fatal("didn't pass output file name")
	// 	// os.Exit(1)
	// 	if fInfo.IsDir() {
	// 		out = filepath.Join(fname, p.PkgName+"_easyjson.go")
	// 	} else {
	// 		if s := strings.TrimSuffix(fname, ".go"); s == fname {
	// 			return errors.New("Filename must end in '.go'")
	// 		} else {
	// 			outName = s + "_easyjson.go"
	// 		}
	// 	}
	// }
	if ok := strings.HasSuffix(out, ".go"); !ok {
		out = out + ".go"
	}

	WriteString( /*subpkg, out, */ pkgCnt /*, apn*/)
}

func WriteString( /*dir, name, */ pkgCnt string /*, apn bool*/) {

	if _, err := os.Stat(subpkg); os.IsNotExist(err) {
		fmt.Println(subpkg, "not exists. Trying to make directory")
		if err := os.Mkdir(subpkg, os.ModePerm); err != nil {
			fmt.Println("Couldn't make directory. Got an error:", err)
			os.Exit(1)
		}
	}
	f = NewFile(pkgCnt)
	defer func() {
		if err := f.Save(filepath.Join(subpkg, out)); err != nil {
			fmt.Println("Couldn't save file. Got an error:", err)
			os.Exit(1)
		}
	}()

	f.Comment("This code is generated by 'csvgen'")
	f.Comment("Do not edit")
	f.Line()

	if err := p.Parse(fname, fInfo.IsDir()); err != nil {
		return
	}
	GenerateCode()
	fmt.Println(p.Error)
	fmt.Println(p.StructMap)
	fmt.Printf("%#v", f)
}

func GenerateCode() {

	for _, v := range p.Structs {
		GenerateFuncs(v)
		f.Line()
	}

}

func GenerateFuncs(vstr parser.StructInfo /*, fields map[string]string*/) {

	// func (this *Type) UnmarshalCSV(in []string) error {
	//  i := 0
	//  if x, err := strconv.ParseBool(in[i]); err == nil {
	//      this.b = x
	//  } else {
	//      return err
	//  }
	//  i++
	//  if x, err := strconv.ParseInt(in[i], 10, 64); err == nil {
	//      this.i = x
	//  } else {
	//      return err
	//  }
	// }
	//
	// // func (this Type) MarshalCSV() (string, error) {
	//      ...marshal logic
	// }

	var unmarshallBody []Code
	var marshallBody []Code

	unmarshallBody = append(unmarshallBody, Id("i").Op(":=").Lit(0))
	marshallBody = append(marshallBody, Id("out").Op(":=").Lit(""))

	for ik, istr := range vstr.Fields {
		var g, s *Statement
		star := ""
		ttype := istr.Type
		if istr.Type[0] == '*' {
			star = "*"
			ttype = istr.Type[1:]
		}
		if star == "*" {
			s = If(
				Id("this").Op(".").Id(istr.Name).Op("==").Id("nil"),
			).Block(
				Return().Qual("errors", "New").Call(Lit("nil pointer found at " + istr.Name + " " + istr.Type)),
			)
		} else {
			s = Null()
		}

		switch ttype {
		case "bool":
			g = If(
				List(Id("x"), Err()).Op(":=").Qual("strconv", "ParseBool").Call(Id("in").Index(Id("i"))),
				Err().Op("!=").Nil(),
			).Add(genReturn(s, star, istr.Name, ttype))
		case "float32":
			fallthrough
		case "float64":
			g = If(
				List(Id("x"), Err()).Op(":=").Qual("strconv", "ParseFloat").Call(List(Id("in").Index(Id("i")), Id(ttype[5:]))),
				Err().Op("!=").Nil(),
			).Add(genReturn(s, star, istr.Name, ttype))
		case "int":
			fallthrough
		case "int8":
			fallthrough
		case "int16":
			fallthrough
		case "int32":
			fallthrough
		case "int64":
			bn := ttype[3:]
			if bn == "" {
				bn = "0"
			}
			g = If(
				List(Id("x"), Err()).Op(":=").Qual("strconv", "ParseInt").Call(List(Id("in").Index(Id("i")), Lit(10), Id(bn))),
				Err().Op("!=").Nil(),
			).Add(genReturn(s, star, istr.Name, ttype))
		case "uint":
			fallthrough
		case "uint8":
			fallthrough
		case "uint16":
			fallthrough
		case "uint32":
			fallthrough
		case "uint64":
			bn := ttype[4:]
			if bn == "" {
				bn = "0"
			}
			g = If(
				List(Id("x"), Err()).Op(":=").Qual("strconv", "ParseUint").
					Call(List(Id("in").Index(Id("i")), Lit(10), Id(bn))),
				Err().Op("!=").Nil(),
			).Add(genReturn(s, star, istr.Name, ttype))
		case "string":
			unmarshallBody = append(unmarshallBody, s)
			g = Op(star).Id("this").Op(".").Id(istr.Name).Op("=").Id("in").Index(Id("i"))
		default:
			g = If(
				Err().Op(":=").Id("this").Op(".").Id(istr.Name).Op(".").Id("UnmarshallCSV").
					Call(Id("in").Index(Id("i"))),
				Err().Op("!=").Nil(),
			).Block(
				Return().Err(),
			)
		}
		unmarshallBody = append(unmarshallBody, g)
		if ik != (len(vstr.Fields) - 1) {
			unmarshallBody = append(unmarshallBody, Id("i").Op("+=").Lit(1))
		} else {
			unmarshallBody = append(unmarshallBody, Return().Id("nil"))
			marshallBody = append(marshallBody, Return().List(Id("out"), Id("nil")))
		}
	}

	f.Comment(vstr.Name + " Unmarshaller func")
	f.Func().Params(
		Id("this").Op("*").Id(vstr.Name),
	).Id("UnmarshalCSV").Params(
		Id("in").Index().String(),
	).Id("error").Block(
		unmarshallBody...,
	)

	f.Line()

	f.Comment(vstr.Name + " Marshaller func")
	f.Func().Params(
		Id("this").Id(vstr.Name),
	).Id("MarshalCSV").Params().
		Parens(Id("string").Op(",").Id("error")).Block(
		marshallBody...,
	)
}

func genReturn(starStmt Code, star string, fieldName string, fieldType string) *Statement {
	var conv *Statement

	//  Check for float and int
	if fieldType[len(fieldType)-2:] != "64" {
		conv = Id(fieldType).Call(Id("x"))
	} else {
		conv = Id("x")
	}
	return Block(
		starStmt,
		Op(star).Id("this").Op(".").Id(fieldName).Op("=").Add(conv),
	).Else().Block(
		Return().Err(),
	)
}
