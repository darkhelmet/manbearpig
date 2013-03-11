package mutations

import (
    "go/ast"
)

type BasicLitFinder struct {
    Value string
    Exps  []*ast.BasicLit
}

func (v *BasicLitFinder) Visit(node ast.Node) ast.Visitor {
    if exp, ok := node.(*ast.BasicLit); ok {
        if exp.Value == v.Value {
            v.Exps = append(v.Exps, exp)
        }
    }
    return v
}

func (v BasicLitFinder) Len() int {
    return len(v.Exps)
}

func (v *BasicLitFinder) Reset() {
    v.Exps = nil
}
