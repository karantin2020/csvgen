package main

import (
    "fmt"
    "log"
    "os"
    "path/filepath"
   . "github.com/dave/jennifer/jen"
   "github.com/karantin2020/cli"
)

var (
	pkg string
	subpkg string
	out string
	apn bool
)

func main() {
	flags := cli.New("Test cli app", "0.0.1")
	flags.StringVarP(&pkg, "pkg", "p", "", "output package")
	flags.StringVarP(&subpkg, "subpkg", "s", "", "output subpkg name")
	flags.StringVarP(&out, "out", "o", "", "output file")
	flags.BoolVarP(&apn, "append", "a", false, "append result of template render to file or not")

	flags.Parse()

	pkgCnt := "package main\n\n"
    if pkg != "" {
    	pkgCnt = "package " + pkg + "\n\n"
    }
	if subpkg != "" {
		pkgCnt = "package " + subpkg + "\n\n"
		pkg = ""
    }
    if out == "" {
		log.Fatal("didn't pass output file name")
		os.Exit(1)
	}

    // path := filepath.Join(subpkg, out + ".go")

	s := ""
	if !apn {
		s = `// This code is generated by 'csvgen'
// DO NOT EDIT

`
	}
	if !apn {
		s = s + pkgCnt
	}

	ls := []Code{}
	ls = append(ls, Qual("fmt", "Println").Call(Lit("Hello, world")))
	f := Func().Id("main").Params().Block( ls... )
	s = s + fmt.Sprintf("%#v", f)
	WriteString(subpkg, out, s, apn)
	// kj := fmt.Sprintf("%#v", f)
	// fmt.Println(kj)
}

func WriteString(dir, name, s string, apn bool) {
	var fileFlag int
	if apn {
		fileFlag = os.O_CREATE|os.O_WRONLY|os.O_APPEND
	} else {
		fileFlag = os.O_CREATE|os.O_TRUNC|os.O_WRONLY
	}

	if _, err := os.Stat(dir); os.IsNotExist(err) {
		os.Mkdir(dir, os.ModePerm)
	}

	path := filepath.Join(dir, name + ".go")

	if f,err := os.OpenFile(path, fileFlag, 0666); err != nil {
		log.Fatal(err)
	} else {
		defer f.Close()
		if _, err:= f.WriteString(s); err != nil {
			log.Fatal(err)
		}
	}
}