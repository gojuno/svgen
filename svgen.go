package main

import (
	"flag"
	"fmt"
	"go/ast"
	"go/types"
	"os"
	"path/filepath"

	"golang.org/x/tools/go/loader"

	"github.com/gojuno/generator"
)

type (
	options struct {
		Package string
	}

	typeValue struct {
		name  string
		value ast.Expr
	}

	typeInfo struct {
		Kind       string //basic type
		FoundScan  bool
		FoundValue bool
		Values     []string
	}

	visitor struct {
		gen         *generator.Generator
		packageInfo *loader.PackageInfo

		types map[string]*typeInfo
	}
)

func main() {
	opts := processFlags()

	packagePath, err := generator.PackageAbsPath(opts.Package)
	if err != nil {
		die(err)
	}

	cfg := loader.Config{
		AllowErrors:         true,
		TypeCheckFuncBodies: func(string) bool { return false },
		TypeChecker: types.Config{
			IgnoreFuncBodies:         true,
			FakeImportC:              true,
			DisableUnusedImportCheck: true,
			Error: func(err error) {},
		},
	}
	cfg.Import(opts.Package)

	outputFile := filepath.Join(packagePath, "scanners_valuers.go")

	if err := os.Remove(outputFile); err != nil && !os.IsNotExist(err) {
		die(err)
	}

	prog, err := cfg.Load()
	if err != nil {
		die(err)
	}

	gen := generator.New(prog)
	gen.SetPackageName(prog.Package(opts.Package).Pkg.Name())
	gen.SetHeader("DO NOT EDIT! This code was generated automatically.")

	v := &visitor{
		gen:         gen,
		packageInfo: prog.Package(opts.Package),
		types:       make(map[string]*typeInfo),
	}

	for _, file := range prog.Package(opts.Package).Files {
		ast.Walk(v, file)
	}

	if err := gen.ProcessTemplate("interface", template, v.types); err != nil {
		die(err)
	}

	if err := gen.WriteToFilename(outputFile); err != nil {
		die(err)
	}
}

func (v *visitor) Visit(node ast.Node) ast.Visitor {
	ts, ok := node.(*ast.ValueSpec)
	if !ok {
		return v
	}

	if _, ok := ts.Type.(*ast.SelectorExpr); ok {
		return nil //type of the constant defined in different package
	}

	ident, ok := ts.Type.(*ast.Ident)
	if !ok {
		return nil
	}

	typeName := ident.Name

	ti, exist := v.types[typeName]
	if !exist {
		ti = &typeInfo{}
	}

	for _, name := range ts.Names {
		ti.Values = append(ti.Values, name.Name)
	}

	if exist {
		return nil
	}

	obj, ok := v.packageInfo.Defs[ts.Names[0]]
	if !ok {
		return nil
	}
	if _, ok := obj.(*types.Const); !ok { //not a constant
		return nil
	}

	nt, ok := obj.Type().(*types.Named)
	if !ok { //not a named type
		return nil
	}

	b, ok := nt.Underlying().(*types.Basic)
	if !ok { //not an alias for a basic type
		return nil
	}

	switch b.Info() {
	case types.IsInteger:
		ti.Kind = "int64"
	case types.IsString:
		ti.Kind = "string"
	default:
		return nil
	}

	for i := 0; i < nt.NumMethods() && (!ti.FoundScan || !ti.FoundValue); i++ {
		method := nt.Method(i)
		switch method.Name() {
		case "Scan":
			ti.FoundScan = true
		case "Value":
			ti.FoundValue = true
		}
	}

	if ti.FoundScan && ti.FoundValue {
		return nil
	}

	v.types[typeName] = ti

	return v
}

const template = `
	{{range $typeName, $typeInfo := . }}
		{{if not $typeInfo.FoundScan }}
			func (t *{{$typeName}}) Scan(i interface{}) error {
				var vv {{$typeName}}
				switch v := i.(type) {
				case nil:
					return nil
				{{if eq $typeInfo.Kind "string"}}case []byte:
					vv = {{$typeName}}(v){{end}}
				case {{$typeInfo.Kind}}:
					vv = {{$typeName}}(v)
				default:
					return fmt.Errorf("can't scan %T into %T", v, t)
				}

				switch vv {
				{{range $value := $typeInfo.Values}}case {{$value}}:{{end}}
				default:
					return fmt.Errorf("invalid value of type {{$typeName}}: %v", *t)
				}

				*t = vv

				return nil
			}
		{{end}}

		{{if not $typeInfo.FoundValue }}
			func (t {{$typeName}}) Value() (driver.Value, error) {
				{{if eq $typeInfo.Kind "string"}}if t == "" {
						return nil, nil
					}
				{{end}}switch t {
				{{range $value := $typeInfo.Values}}case {{$value}}:{{end}}
				default:
					return nil, fmt.Errorf("invalid value of type {{$typeName}}: %v", t)
				}

				return {{$typeInfo.Kind}}(t), nil
			}
		{{end}}
	{{end}}`

func processFlags() *options {
	var (
		input = flag.String("i", "", "import path of the package containing type declarations")
	)

	flag.Parse()

	if *input == "" {
		flag.Usage()
		os.Exit(1)
	}

	return &options{
		Package: *input,
	}
}

func die(err error) {
	fmt.Fprintf(os.Stderr, "%v\n", err)
	os.Exit(1)
}
