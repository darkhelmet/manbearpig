package mutations

import (
    "go/ast"
    "go/token"
)

type BinaryExpressionFinder struct {
    Token token.Token
    Exps  []*ast.BinaryExpr
}

func (v *BinaryExpressionFinder) Visit(node ast.Node) ast.Visitor {
    if exp, ok := node.(*ast.BinaryExpr); ok {
        if exp.Op == v.Token {
            v.Exps = append(v.Exps, exp)
        }
    }
    return v
}

func (v BinaryExpressionFinder) Len() int {
    return len(v.Exps)
}

func (v *BinaryExpressionFinder) Reset() {
    v.Exps = nil
}
