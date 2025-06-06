package main

import (
	"fmt"
	"go/ast"
	"strconv"

	"github.com/dave/jennifer/jen"
)

type ParamRenamer func(index int, name string) string

var (
	RenamePostfix = ParamRenamer(func(_ int, name string) string {
		return name + "_"
	})
	RenameUnnamed = ParamRenamer(func(index int, name string) string {
		if name != "" {
			return name
		}
		return fmt.Sprintf("arg%d", index)
	})
)

func (rename ParamRenamer) rename(index int, p *Param) {
	if rename != nil {
		p.Name = rename(index, p.Name)
	}
}

type Param struct {
	Name       string
	Type       ast.Expr
	TypeJen    jen.Code
	IsVariadic bool
}

type Params []Param

func (params Params) ToSignatureParams(renamers ...ParamRenamer) []jen.Code {
	out := make([]jen.Code, 0, len(params))
	for i, param := range params {
		for _, r := range renamers {
			r.rename(i, &param)
		}
		out = append(out, jen.Id(param.Name).Add(param.TypeJen))
	}
	return out
}

func (params Params) ToCallParams(useVariadic bool, renamers ...ParamRenamer) []jen.Code {
	out := make([]jen.Code, 0, len(params))
	for i, param := range params {
		for _, r := range renamers {
			r.rename(i, &param)
		}
		t := jen.Id(param.Name)
		if param.IsVariadic && useVariadic {
			t = t.Add(jen.Op("..."))
		}
		out = append(out, t)
	}
	return out
}

func (params Params) ToTypeIds() []jen.Code {
	out := make([]jen.Code, 0, len(params))
	for _, param := range params {
		out = append(out, jen.Id(param.Name))
	}
	return out
}

func OptionalTypes(codes []jen.Code) jen.Code {
	if len(codes) != 0 {
		return jen.Types(codes...)
	}
	return nil
}

func paramsFromFieldList(fl *ast.FieldList, imports importMap) Params {
	if fl == nil {
		return nil
	}
	codes := make(Params, 0, fl.NumFields())
	for _, l := range fl.List {
		typ := exprToJen(l.Type, imports)
		_, isVariadic := l.Type.(*ast.Ellipsis)
		if len(l.Names) == 0 {
			codes = append(codes, Param{Type: l.Type, TypeJen: typ, IsVariadic: isVariadic})
			continue
		}
		for _, n := range l.Names {
			codes = append(codes, Param{Name: n.Name, Type: l.Type, TypeJen: typ, IsVariadic: isVariadic})
		}
	}
	return codes
}

type importMap map[string]string

func exprToJen(e ast.Expr, imports importMap) jen.Code {
	switch e := e.(type) {
	case *ast.Ellipsis:
		return jen.Op("...").Add(exprToJen(e.Elt, imports))

	case *ast.ArrayType:
		if e.Len == nil {
			return jen.Index().Add(exprToJen(e.Elt, imports))
		}
		l, err := strconv.ParseInt(e.Len.(*ast.BasicLit).Value, 10, 64)
		if err != nil {
			panic("invalid lit")
		}
		return jen.Index(jen.Lit(l)).Add(exprToJen(e.Elt, imports))

	case *ast.StructType:
		return jen.StructFunc(func(g *jen.Group) {
			if e.Fields == nil {
				return
			}
			for _, f := range e.Fields.List {
				stmt := g.Null()
				for i, n := range f.Names {
					if i > 0 {
						stmt.Op(",")
					}
					stmt.Id(n.Name)
				}
				stmt.Add(exprToJen(f.Type, imports))
				// NOTE: struct tags are currently not supported
			}
		})

	case *ast.StarExpr:
		return jen.Op("*").Add(exprToJen(e.X, imports))

	case *ast.FuncType:
		return jen.Func().
			Params(paramsFromFieldList(e.Params, imports).ToSignatureParams()...).
			Params(paramsFromFieldList(e.Results, imports).ToSignatureParams()...)

	case *ast.InterfaceType:
		return jen.InterfaceFunc(func(g *jen.Group) {
			for _, m := range e.Methods.List {
				if len(m.Names) == 0 {
					panic("embedded interfaces are not supported")
				}
				fn := m.Type.(*ast.FuncType)
				g.Id(m.Names[0].Name).
					Params(paramsFromFieldList(fn.Params, imports).ToSignatureParams()...).
					Params(paramsFromFieldList(fn.Results, imports).ToSignatureParams()...)
			}
		})

	case *ast.MapType:
		return jen.Map(exprToJen(e.Key, imports)).Add(exprToJen(e.Value, imports))

	case *ast.ChanType:
		stmt := jen.Chan()
		if e.Dir == ast.RECV {
			stmt = jen.Op("<-").Add(stmt)
		}
		if e.Dir == ast.SEND {
			stmt.Op("<-")
		}
		return stmt.Add(exprToJen(e.Value, imports))

	case *ast.Ident:
		return jen.Id(e.Name)

	case *ast.SelectorExpr:
		if e.X == nil {
			return jen.Id(e.Sel.Name)
		}
		x := e.X.(*ast.Ident).Name
		if r, ok := imports[x]; ok {
			x = r
		}
		return jen.Qual(x, e.Sel.Name)

	case *ast.UnaryExpr:
		return jen.Op(e.Op.String()).Add(exprToJen(e.X, imports))

	case *ast.IndexExpr:
		return jen.Add(exprToJen(e.X, imports)).Index(exprToJen(e.Index, imports))

	default:
		panic(fmt.Errorf("unsupported ast.Expr: %T", e))
	}
}
