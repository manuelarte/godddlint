package astutils

import (
	"go/ast"
	"slices"
)

// GetRcvName gets the name of an Expr.
func GetRcvName(expr ast.Expr) (string, bool) {
	switch expr := expr.(type) {
	case *ast.Ident:
		return expr.Name, true
	case *ast.StarExpr:
		return GetRcvName(expr.X)
	default:
		return "", false
	}
}

func GetFieldsThatMustDefensiveCopy(st *ast.StructType) []*ast.Field {
	fields := make([]*ast.Field, 0)

	for _, field := range st.Fields.List {
		_, isMap := field.Type.(*ast.MapType)

		_, isSlice := field.Type.(*ast.ArrayType)
		if isMap || isSlice {
			fields = append(fields, field)
		}
	}

	return fields
}

func IsPotentialConstructor(f *ast.FuncDecl) bool {
	if f.Recv != nil {
		return false
	}

	results := f.Type.Results.List
	if len(results) == 1 {
		return true
	}

	if len(results) == 2 {
		if ident, ok := results[1].Type.(*ast.Ident); ok {
			if ident.Name == "error" {
				return true
			}
		}
	}

	return false
}

func IsPrimitiveType(f *ast.Field) bool {
	ident, ok := f.Type.(*ast.Ident)
	if !ok {
		return false
	}

	primitiveTypes := []string{
		"bool", "byte", "complex64", "complex128",
		"float32", "float64", "int", "int8", "int16", "int32", "int64", "string",
	}

	return slices.Contains(primitiveTypes, ident.Name)
}

func IsErrorsNewFun(callExpr *ast.CallExpr) bool {
	if se, ok := callExpr.Fun.(*ast.SelectorExpr); ok {
		ident, isIdent := se.X.(*ast.Ident)
		if !isIdent {
			return false
		}

		return ident.Name == "errors" && se.Sel.Name == "New"
	}

	return false
}

// FuncResultError returns the position in the function results where the error is.
// e.g. func myFunction() (int, error), would return (1, true).
func FuncResultError(f *ast.FuncDecl) (int, bool) {
	if f.Type.Results == nil {
		return -1, false
	}

	for i, result := range f.Type.Results.List {
		ident, ok := result.Type.(*ast.Ident)
		if !ok {
			continue
		}

		if ident.Name == "error" {
			return i, true
		}
	}

	return -1, false
}
