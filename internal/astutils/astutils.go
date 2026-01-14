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
