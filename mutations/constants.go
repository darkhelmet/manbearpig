package mutations

import (
    "go/ast"
    "log"
    "path/filepath"
)

func init() {
    Mutations["0"] = NewConstantMutation("0", "1")
}

type ConstantMutation struct {
    commonMutation
    exps      []*ast.BinaryExpr
    expFinder *BasicLitFinder
    mutations []string
}

func (bop *ConstantMutation) Prepare(src string, logger *log.Logger) {
    bop.src = src
    bop.logger = logger
    bop.expFinder.Reset()
    bop.parse(bop.expFinder)
}

func (bop *ConstantMutation) Run() {
    bop.logf("found %d occurrence(s) of %s in %s", bop.expFinder.Len(), bop.expFinder.Value, filepath.Base(bop.src))
    for _, m := range bop.mutations {
        if bop.expFinder.Len() > 0 {
            bop.logf("mutating %s to %s", bop.expFinder.Value, m)
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

func (bop *ConstantMutation) perform(n int, exp *ast.BasicLit, value string) error {
    exp.Value = value
    defer func() {
        exp.Value = bop.expFinder.Value
    }()
    return bop.runTests(n)
}

func NewConstantMutation(value string, mutations ...string) *ConstantMutation {
    return &ConstantMutation{
        mutations: mutations,
        expFinder: &BasicLitFinder{Value: value},
    }
}
