package main

import (
	"reflect"
	"testing"

	. "github.com/dave/jennifer/jen"
	"github.com/karantin2020/csvgen/parser"
)

//go:generate ./csvgen -p tests -s tests -f tests/fixture/test.go -o test
//go:generate ./csvgen -f tests/fixture

func Test_main(t *testing.T) {
	tests := []struct {
		name string
	}{
	// TODO: Add test cases.
	}
	for range tests {
		t.Run(tt.name, func(t *testing.T) {
			main()
		})
	}
}

func TestWriteString(t *testing.T) {
	type args struct {
		pkgCnt string
	}
	tests := []struct {
		name string
		args args
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			WriteString(tt.args.pkgCnt)
		})
	}
}

func TestGenerateCode(t *testing.T) {
	tests := []struct {
		name string
	}{
	// TODO: Add test cases.
	}
	for range tests {
		t.Run(tt.name, func(t *testing.T) {
			GenerateCode()
		})
	}
}

func TestGenerateFuncs(t *testing.T) {
	type args struct {
		vstr parser.StructInfo
	}
	tests := []struct {
		name string
		args args
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			GenerateFuncs(tt.args.vstr)
		})
	}
}

func Test_genReturn(t *testing.T) {
	type args struct {
		star      string
		fieldName string
		fieldType string
	}
	tests := []struct {
		name string
		args args
		want *Statement
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := genReturn(tt.args.star, tt.args.fieldName, tt.args.fieldType); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("genReturn() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_marshalBody(t *testing.T) {
	type args struct {
		typeRes *Statement
	}
	tests := []struct {
		name string
		args args
		want *Statement
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := marshalBody(tt.args.typeRes); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("marshalBody() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_nilCheck(t *testing.T) {
	type args struct {
		star     string
		iname    string
		itype    string
		marshall bool
	}
	tests := []struct {
		name string
		args args
		want *Statement
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := nilCheck(tt.args.star, tt.args.iname, tt.args.itype, tt.args.marshall); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("nilCheck() = %v, want %v", got, tt.want)
			}
		})
	}
}
