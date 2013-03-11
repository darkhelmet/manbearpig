package mutations

import (
    "go/ast"
    "go/token"
    "log"
    "path/filepath"
)

func init() {
    Mutations["=="] = NewBinaryOperatorMutation(token.EQL, token.NEQ)
    Mutations["!="] = NewBinaryOperatorMutation(token.NEQ, token.EQL)

    Mutations[">"] = NewBinaryOperatorMutation(token.GTR, token.LSS, token.GEQ, token.LEQ)
    Mutations["<"] = NewBinaryOperatorMutation(token.LSS, token.GTR, token.LEQ, token.GEQ)

    Mutations[">="] = NewBinaryOperatorMutation(token.GEQ, token.GTR, token.LEQ, token.LSS)
    Mutations["<="] = NewBinaryOperatorMutation(token.LEQ, token.LSS, token.GEQ, token.GTR)

    Mutations["||"] = NewBinaryOperatorMutation(token.LOR, token.LAND)
    Mutations["&&"] = NewBinaryOperatorMutation(token.LAND, token.LOR)

    Mutations["|"] = NewBinaryOperatorMutation(token.OR, token.AND)
    Mutations["&"] = NewBinaryOperatorMutation(token.AND, token.OR)
}

type BinaryOperatorMutation struct {
    commonMutation
    exps      []*ast.BinaryExpr
    expFinder *BinaryExpressionFinder
    mutations []token.Token
}

func (bop *BinaryOperatorMutation) Prepare(src string, logger *log.Logger) {
    bop.src = src
    bop.logger = logger
    bop.expFinder.Reset()
    bop.parse(bop.expFinder)
}

func (bop *BinaryOperatorMutation) Run() {
    bop.logf("found %d occurrence(s) of %s in %s", bop.expFinder.Len(), bop.expFinder.Token, filepath.Base(bop.src))
    for _, m := range bop.mutations {
        if bop.expFinder.Len() > 0 {
            bop.logf("mutating %s to %s", bop.expFinder.Token, m)
            for index, exp := range bop.expFinder.Exps {
                err := bop.perform(index+1, exp, m)
                if err != nil {
                    bop.logf("mutation failed!")
                }
            }
        }
    }
    err := bop.printFile()
    if err != nil {
        bop.fatalf("failed to restore %s to original state, further mutations tainted: %s", err)
    }
}

func (bop *BinaryOperatorMutation) perform(n int, exp *ast.BinaryExpr, tok token.Token) error {
    exp.Op = tok
    defer func() {
        exp.Op = bop.expFinder.Token
    }()
    return bop.runTests(n)
}

func NewBinaryOperatorMutation(tok token.Token, mutations ...token.Token) *BinaryOperatorMutation {
    return &BinaryOperatorMutation{
        mutations: mutations,
        expFinder: &BinaryExpressionFinder{Token: tok},
    }
}
