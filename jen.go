package main

import (
	"fmt"
	"go/types"

	"github.com/dave/jennifer/jen"
)

type ParamRenamer func(index int, name string) string

var (
	RenamePostfix = ParamRenamer(func(_ int, name string) string {
		return name + "_"
	})
	RenameUnnamed = ParamRenamer(func(index int, name string) string {
		return fmt.Sprintf("arg%d", index)
	})
)

func (rename ParamRenamer) rename(index int, p *Param) {
	if rename != nil {
		p.Name = rename(index, p.Name)
	}
}

type Param struct {
	Name         string
	Type         jen.Code
	VariadicType jen.Code
	AstType      types.Type
}

type Params []Param

func (params Params) ToSignatureParams(renamer ParamRenamer) []jen.Code {
	out := make([]jen.Code, 0, len(params))
	for i, param := range params {
		renamer.rename(i, &param)
		// either VariadicType or Type is set so both can safely be appended
		out = append(out, jen.Id(param.Name).Add(param.VariadicType, param.Type))
	}
	return out
}

func (params Params) ToCallParams(renamer ParamRenamer) []jen.Code {
	out := make([]jen.Code, 0, len(params))
	for i, param := range params {
		renamer.rename(i, &param)
		out = append(out, jen.Id(param.Name))
	}
	return out
}

func signatureToJen(sig *types.Signature) jen.Code {
	return jen.
		Params(tupleToJen(sig.Params(), sig.Variadic()).ToSignatureParams(nil)...).
		Params(tupleToJen(sig.Results(), false).ToSignatureParams(nil)...)
}

func tupleToJen(tuple *types.Tuple, variadic bool) Params {
	codes := make(Params, 0, tuple.Len())

	for i := 0; i < tuple.Len(); i++ {
		v := tuple.At(i)
		param := Param{Name: v.Name(), AstType: v.Type()}

		if i == tuple.Len()-1 && variadic {
			slice := v.Type().(*types.Slice)
			param.VariadicType = jen.Op("...").Add(typeToJen(slice.Elem()))
		} else {
			param.Type = typeToJen(v.Type())
		}

		codes = append(codes, param)
	}

	return codes
}

func typeToJen(v types.Type) jen.Code {
	switch v := v.(type) {
	case *types.Basic:
		return jen.Id(v.Name())

	case *types.Array:
		return jen.Index(jen.Lit(int(v.Len()))).Add(typeToJen(v.Elem()))

	case *types.Slice:
		return jen.Index().Add(typeToJen(v.Elem()))

	case *types.Struct:
		return jen.StructFunc(func(g *jen.Group) {
			for i := 0; i < v.NumFields(); i++ {
				field := v.Field(i)
				if field.Embedded() {
					g.Add(typeToJen(field.Type()))
				} else {
					g.Id(field.Name()).Add(typeToJen(field.Type()))
				}
			}
		})

	case *types.Pointer:
		return jen.Op("*").Add(typeToJen(v.Elem()))

	case *types.Signature:
		return jen.Func().Add(signatureToJen(v))

	case *types.Interface:
		return jen.InterfaceFunc(func(g *jen.Group) {
			for i := 0; i < v.NumEmbeddeds(); i++ {
				g.Add(typeToJen(v.EmbeddedType(i)))
			}
			for i := 0; i < v.NumExplicitMethods(); i++ {
				fun := v.ExplicitMethod(i)
				sig := fun.Type().(*types.Signature)
				g.Id(fun.Name()).Add(signatureToJen(sig))
			}
		})

	case *types.Map:
		return jen.Map(typeToJen(v.Key())).Add(typeToJen(v.Elem()))

	case *types.Chan:
		stmt := jen.Chan()
		if v.Dir() == types.RecvOnly {
			stmt = jen.Op("<-").Add(stmt)
		}
		if v.Dir() == types.SendOnly {
			stmt.Op("<-")
		}
		return stmt.Add(typeToJen(v.Elem()))

	case *types.Named:
		obj := v.Obj()
		if obj.Pkg() == nil {
			return jen.Id(obj.Name())
		}
		return jen.Qual(obj.Pkg().Path(), obj.Name())

	default:
		panic(fmt.Errorf("unsupported *ast.Var type: %T", v))
	}
}

func isNillable(typ types.Type) bool {
	switch t := typ.(type) {
	case *types.Array, *types.Slice, *types.Pointer, *types.Signature, *types.Interface, *types.Map, *types.Chan:
		return true
	case *types.Named:
		return isNillable(t.Underlying())
	}
	return false
}
