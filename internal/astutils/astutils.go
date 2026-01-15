package astutils

import (
	"go/ast"
	"slices"
)

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

// GetFieldsThatMustDefensiveCopy get the fields in the struct declaration that needs to be defensive copied.
// so far map and slices.
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

// IsErrorsNewFun checks if the call expression is an errors.New function.
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

// IsPrimitiveType checks if the field is a primitive type.
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

// IsPotentialConstructor checks if the function is a potential constructor.
// if return either one element or two being the second an error.
func IsPotentialConstructor(f *ast.FuncDecl) bool {
	if f.Recv != nil {
		return false
	}

	if f.Type.Results == nil {
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
