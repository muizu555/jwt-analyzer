package analyzer

import (
	"go/ast"
	"go/token"

	"golang.org/x/tools/go/analysis"
)

var IssuerAnalyzer = &analysis.Analyzer{
	Name: "jwtissuercheck",
	Doc:  "checks that JWT issuer is explicitly validated after parsing",
	Run:  runIssuerCheck,
}

func runIssuerCheck(pass *analysis.Pass) (any, error) {
	for _, file := range pass.Files {
		ast.Inspect(file, func(n ast.Node) bool {
			fn, ok := n.(*ast.FuncDecl)
			if !ok {
				return true
			}

			// 各関数の中で JWT パース呼び出しがあるかと Issuer チェックがあるかを調べる
			hasParse := false
			hasIssuerCheck := false

			ast.Inspect(fn.Body, func(n ast.Node) bool {
				switch x := n.(type) {
				case *ast.CallExpr:
					if getFuncName(x.Fun) == "ParseWithClaims" {
						hasParse = true
					}
				case *ast.BinaryExpr:
					if x.Op == token.EQL || x.Op == token.NEQ {
						// 左辺にIssuerが含まれているか確認
						if isIssuerCheck(x.X) {
							hasIssuerCheck = true
						}

						// 右辺にIssuerが含まれているか確認
						if isIssuerCheck(x.Y) {
							hasIssuerCheck = true
						}
					}
				}
				return true
			})

			// JWTパースがあるのにIssuerチェックが無ければ警告
			if hasParse && !hasIssuerCheck {
				pass.Reportf(fn.Pos(), "JWT is parsed but issuer is not checked in function '%s'", fn.Name.Name)
			}

			return true
		})
	}
	return nil, nil
}

// Issuerフィールドの比較を検出するヘルパー関数
func isIssuerCheck(expr ast.Expr) bool {
	// 直接 claims.Issuer の形のチェック
	if sel, ok := expr.(*ast.SelectorExpr); ok {
		return sel.Sel.Name == "Issuer"
	}
	return false
}
