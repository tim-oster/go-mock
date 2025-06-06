package main

import (
	"fmt"
	"go/ast"
	"log"
	"path"
	"strings"
	"unicode"

	"github.com/dave/jennifer/jen"
)

const (
	testingPkg = "testing"
	mockPkg    = "github.com/stretchr/testify/mock"
)

var (
	mockReceiver = jen.Id("x")
)

type File struct {
	imports    []*ast.ImportSpec
	pkgName    string
	filename   string
	interfaces []Interface
}

func parseFiles(f *ast.File, filename string, targets map[string]string) []File {
	var files []File
	ff := File{
		imports:  f.Imports,
		filename: filename,
		pkgName:  f.Name.Name,
	}
	ast.Inspect(f, ff.genDecl)
	if !ff.filterInterfaces(targets) {
		return nil
	}
	files = append(files, ff)
	return files
}

func (f *File) genDecl(node ast.Node) bool {
	spec, ok := node.(*ast.TypeSpec)
	if !ok {
		return true
	}
	iface, ok := spec.Type.(*ast.InterfaceType)
	if !ok {
		return false
	}
	f.interfaces = append(f.interfaces, Interface{
		name: spec.Name.Name,
		spec: spec,
		typ:  iface,
	})
	return false
}

func (f *File) filterInterfaces(targets map[string]string) bool {
	var filtered []Interface

	for _, iface := range f.interfaces {
		rename, ok := targets[iface.name]
		if !ok {
			continue
		}
		if len(rename) > 0 {
			iface.name = rename
		}
		filtered = append(filtered, iface)
	}

	f.interfaces = filtered
	return len(f.interfaces) > 0
}

func (f *File) generate() {
	imports := importMap{}
	for _, spec := range f.imports {
		pathValue := strings.Trim(spec.Path.Value, "\"")
		if spec.Name != nil {
			name := strings.Trim(spec.Name.Name, "\"")
			imports[name] = pathValue
			continue
		}
		pkgName := path.Base(pathValue)
		if pkgName != pathValue {
			imports[pkgName] = pathValue
		}
	}

	for _, iface := range f.interfaces {
		g := jen.NewFile(f.pkgName)
		g.PackageComment("Code generated by " + binaryName + ". DO NOT EDIT.")
		g.PackageComment("version: " + binaryName + " " + binaryVersion)
		g.PackageComment("source: " + f.filename)

		actualName := iface.spec.Name.Name
		generatedName := iface.name
		if actualName != generatedName {
			actualName += " (renamed to: " + generatedName + ")"
		}
		g.PackageComment("interface: " + actualName)

		g.PackageComment(fmt.Sprintf("flags: keepctx=%t unexported=%t", *keepCtx, *unexported))

		g.Line()

		iface.generate(g, imports)

		filename := path.Base(f.filename)
		filename = strings.TrimSuffix(filename, ".go")
		filename += "_mock_" + strings.ToLower(iface.name)
		filename += ".gen.go"
		err := g.Save(filename)
		if err != nil {
			log.Fatal(err)
		}
	}
}

type Interface struct {
	name string
	spec *ast.TypeSpec
	typ  *ast.InterfaceType
}

func (iface *Interface) generate(g *jen.File, imports importMap) {
	mockName := "Mock"
	if *unexported {
		mockName = "mock"
	}
	if unicode.IsLower(rune(iface.name[0])) {
		mockName += "_"
	}
	mockName += iface.name

	// Compile time interface implementation assert. Do not generate for generics, as constraints cannot inferred reliably.
	if iface.spec.TypeParams == nil {
		g.Var().Id("_").Id(iface.spec.Name.Name).Op("=").Params(jen.Op("*").Id(mockName)).Params(jen.Nil())
	}

	// mock struct that implements mock.Mock
	genericTypes := paramsFromFieldList(iface.spec.TypeParams, imports)
	g.Type().Id(mockName).Add(OptionalTypes(genericTypes.ToSignatureParams())).StructFunc(func(g *jen.Group) {
		g.Qual(mockPkg, "Mock")
	})

	for i := range iface.typ.Methods.NumFields() {
		m := iface.typ.Methods.List[i]
		fn, ok := m.Type.(*ast.FuncType)
		if !ok {
			continue
		}
		iface.generateForMethod(g, mockName, m.Names[0].Name, fn, imports)
	}
}

func (iface *Interface) generateForMethod(g *jen.File, structName, fnName string, fn *ast.FuncType, imports importMap) {
	mockedInput := paramsFromFieldList(fn.Params, imports)

	if !*keepCtx && len(mockedInput) != 0 {
		p0 := mockedInput[0]
		if sel, ok := p0.Type.(*ast.SelectorExpr); ok {
			packageName := sel.X.(*ast.Ident).Name
			isContextPkg := packageName == "context" || imports[packageName] == "context"
			if isContextPkg && sel.Sel != nil && sel.Sel.Name == "Context" {
				mockedInput = mockedInput[1:]
			}
		}
	}
	m := Method{
		structName:   structName,
		name:         fnName,
		genericTypes: paramsFromFieldList(iface.spec.TypeParams, imports),
		params:       paramsFromFieldList(fn.Params, imports),
		results:      paramsFromFieldList(fn.Results, imports),
		mockedInput:  mockedInput,
	}

	m.generateMockImplementation(g)
	g.Line()
	m.generateCallStruct(g)
	g.Line()
	m.generateOnMethods(g)
	g.Line()
	m.generateAssertMethods(g)
	g.Line()
}

type Method struct {
	structName   string
	name         string
	genericTypes Params
	params       Params
	results      Params
	mockedInput  Params
}

func (m *Method) generateMockImplementation(g *jen.File) {
	optionalGenerics := OptionalTypes(m.genericTypes.ToTypeIds())

	g.Func().
		Params(jen.Add(mockReceiver).Op("*").Id(m.structName).Add(optionalGenerics)).
		Id(m.name).
		Params(m.params.ToSignatureParams(RenameUnnamed, RenamePostfix)...).
		Params(m.results.ToSignatureParams()...).
		BlockFunc(func(g *jen.Group) {
			g.Id("args").Op(":=").Add(mockReceiver).Dot("Called").Call(m.mockedInput.ToCallParams(false, RenameUnnamed, RenamePostfix)...)

			g.If(jen.Len(jen.Id("args")).Op(">").Lit(0)).BlockFunc(func(g *jen.Group) {
				g.If(
					jen.List(jen.Id("t"), jen.Id("ok")).Op(":=").
						Id("args").Dot("Get").Call(jen.Lit(0)).Op(".").Params(jen.Id(m.returnFuncName()).Add(optionalGenerics)),
					jen.Id("ok"),
				).BlockFunc(func(g *jen.Group) {
					invocation := jen.Id("t").Params(m.params.ToCallParams(true, RenameUnnamed, RenamePostfix)...)
					if len(m.results) > 0 {
						g.Return(invocation)
					} else {
						g.Add(invocation)
					}
				})
			})

			if len(m.results) > 0 {
				returnIds := make([]jen.Code, 0, len(m.results))
				for i, result := range m.results {
					id := jen.Id(fmt.Sprintf("r%d", i))
					returnIds = append(returnIds, id)

					g.Var().Add(id).Add(result.TypeJen)
					g.If(
						jen.Id("v").Op(":=").Id("args").Dot("Get").Call(jen.Lit(i)),
						jen.Id("v").Op("!=").Nil(),
					).BlockFunc(func(g *jen.Group) {
						g.Add(id).Op("=").Id("v").Op(".").Params(result.TypeJen)
					})
				}
				g.Return(returnIds...)
			}
		})
}

func (m *Method) callStructName() string {
	structName := []byte(m.structName)
	structName[0] = byte(unicode.ToLower(rune(structName[0])))
	return string(structName) + "_" + m.name
}

func (m *Method) returnFuncName() string {
	return m.callStructName() + "_ReturnFunc"
}

func (m *Method) generateCallStruct(g *jen.File) {
	mockCallStruct := m.callStructName()
	optionalGenerics := OptionalTypes(m.genericTypes.ToTypeIds())

	// actual call struct
	g.Type().Id(mockCallStruct).Add(OptionalTypes(m.genericTypes.ToSignatureParams())).Struct(
		jen.Op("*").Qual(mockPkg, "Call"),
	).Line()

	// generate return function type for mocking custom logic
	g.Type().Id(m.returnFuncName()).Add(OptionalTypes(m.genericTypes.ToSignatureParams())).
		Func().Params(m.params.ToSignatureParams()...).
		Params(m.results.ToSignatureParams()...).
		Line()

	// mockStruct.Return
	g.Func().
		Params(jen.Id("c").Op("*").Id(mockCallStruct).Add(optionalGenerics)).
		Id("Return").
		Params(m.results.ToSignatureParams(RenameUnnamed)...).
		Params(jen.Op("*").Qual(mockPkg, "Call")).
		BlockFunc(func(g *jen.Group) {
			g.Return(jen.Id("c").Dot("Call").Dot("Return").CallFunc(func(g *jen.Group) {
				for i, param := range m.results {
					RenameUnnamed.rename(i, &param)
					g.Id(param.Name)
				}
			}))
		}).
		Line()

	// mockStruct.ReturnFn
	g.Func().
		Params(jen.Id("c").Op("*").Id(mockCallStruct).Add(optionalGenerics)).
		Id("ReturnFn").
		Params(jen.Id("fn").Id(m.returnFuncName()).Add(optionalGenerics)).
		Params(jen.Op("*").Qual(mockPkg, "Call")).
		BlockFunc(func(g *jen.Group) {
			g.Return(jen.Id("c").Dot("Call").Dot("Return").Call(jen.Id("fn")))
		}).
		Line()
}

func (m *Method) generateOnMethods(g *jen.File) {
	optionalGenerics := OptionalTypes(m.genericTypes.ToTypeIds())

	genFunc := func(name string, sigParams []jen.Code, callParams []jen.Code) {
		callParams = append([]jen.Code{jen.Lit(m.name)}, callParams...)

		g.Func().
			Params(jen.Add(mockReceiver).Op("*").Id(m.structName).Add(optionalGenerics)).
			Id(name).
			Params(sigParams...).
			Params(jen.Op("*").Id(m.callStructName()).Add(optionalGenerics)).
			BlockFunc(func(g *jen.Group) {
				g.Return(jen.Op("&").Id(m.callStructName()).Add(optionalGenerics).Values(jen.Dict{
					jen.Id("Call"): jen.Add(mockReceiver).Dot("On").Call(callParams...),
				}))
			}).
			Line()
	}

	// mockStruct.On_xxx
	genFunc(
		"On_"+m.name,
		m.mockedInput.ToSignatureParams(RenameUnnamed, RenamePostfix),
		m.mockedInput.ToCallParams(false, RenameUnnamed, RenamePostfix),
	)

	// mockStruct.On_xxx_Any
	if len(m.mockedInput) > 0 {
		var anythings []jen.Code
		for range m.mockedInput {
			anythings = append(anythings, jen.Qual(mockPkg, "Anything"))
		}
		genFunc("On_"+m.name+"_Any", nil, anythings)
	}

	// mockStruct.On_xxx_Interface
	var interfacedParams Params
	for _, param := range m.mockedInput {
		param.TypeJen = jen.Any()
		interfacedParams = append(interfacedParams, param)
	}
	if len(interfacedParams) > 0 {
		genFunc(
			"On_"+m.name+"_Interface",
			interfacedParams.ToSignatureParams(RenameUnnamed, RenamePostfix),
			interfacedParams.ToCallParams(false, RenameUnnamed, RenamePostfix),
		)
	}
}

func (m *Method) generateAssertMethods(g *jen.File) {
	optionalGenerics := OptionalTypes(m.genericTypes.ToTypeIds())

	genFunc := func(name, implName string, sigParams []jen.Code, callParams []jen.Code) {
		sigParams = append([]jen.Code{jen.Id("t").Op("*").Qual(testingPkg, "T")}, sigParams...)
		callParams = append([]jen.Code{jen.Id("t"), jen.Lit(m.name)}, callParams...)

		g.Func().
			Params(jen.Add(mockReceiver).Op("*").Id(m.structName).Add(optionalGenerics)).
			Id(name).
			Params(sigParams...).
			Params(jen.Bool()).
			BlockFunc(func(g *jen.Group) {
				g.Return(jen.Add(mockReceiver).Dot(implName).Params(callParams...))
			}).
			Line()
	}

	// Assert_xxx_Called function
	genFunc(
		"Assert_"+m.name+"_Called",
		"AssertCalled",
		m.mockedInput.ToSignatureParams(RenameUnnamed, RenamePostfix),
		m.mockedInput.ToCallParams(false, RenameUnnamed, RenamePostfix),
	)

	// Assert_xxx_NumberOfCalls function
	genFunc(
		"Assert_"+m.name+"_NumberOfCalls",
		"AssertNumberOfCalls",
		[]jen.Code{jen.Id("expectedCalls").Int()},
		[]jen.Code{jen.Id("expectedCalls")},
	)

	// Assert_xxx_NotCalled function
	genFunc(
		"Assert_"+m.name+"_NotCalled",
		"AssertNotCalled",
		m.mockedInput.ToSignatureParams(RenameUnnamed, RenamePostfix),
		m.mockedInput.ToCallParams(false, RenameUnnamed, RenamePostfix),
	)
}
