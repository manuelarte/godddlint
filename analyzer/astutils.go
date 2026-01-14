package analyzer

import "go/ast"

// getRcvName gets the name of an Expr.
func getRcvName(expr ast.Expr) (string, bool) {
	switch expr := expr.(type) {
	case *ast.Ident:
		return expr.Name, true
	case *ast.StarExpr:
		return getRcvName(expr.X)
	default:
		return "", false
	}
}

func potentialConstructor(f *ast.FuncDecl) bool {
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
