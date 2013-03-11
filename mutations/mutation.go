package mutations

import (
    "bytes"
    "fmt"
    "go/ast"
    "go/parser"
    "go/printer"
    "go/token"
    "log"
    "os"
    "os/exec"
    "path/filepath"
)

var (
    Mutations = make(map[string]Mutation)
)

type Mutation interface {
    Prepare(string, *log.Logger)
    Run()
    FailureCount() int
}

type commonMutation struct {
    failed int
    src    string
    logger *log.Logger
    fset   *token.FileSet
    file   *ast.File
}

func (cm *commonMutation) FailureCount() int {
    return cm.failed
}

func (cm *commonMutation) parse(v ast.Visitor) {
    cm.fset = token.NewFileSet()
    f, err := parser.ParseFile(cm.fset, cm.src, nil, 0)
    if err != nil {
        cm.logger.Fatalf("failed to parse %s: %s", cm.src, err)
    }
    cm.file = f
    ast.Walk(v, f)
}

func (cm *commonMutation) logf(format string, v ...interface{}) {
    cm.logger.Printf(format, v...)
}

func (cm *commonMutation) fatalf(format string, v ...interface{}) {
    cm.logger.Fatalf(format, v...)
}

func (cm *commonMutation) printFile() error {
    file, err := os.OpenFile(cm.src, os.O_WRONLY|os.O_TRUNC, 0)
    if err != nil {
        return fmt.Errorf("failed to open output file %s: %s", cm.src, err)
    }
    defer file.Close()

    err = printer.Fprint(file, cm.fset, cm.file)
    if err != nil {
        return fmt.Errorf("failed to write AST to file: %s", err)
    }
    return nil
}

func (cm *commonMutation) runTests(n int) error {
    err := cm.printFile()
    if err != nil {
        return err
    }

    cmd := exec.Command("go", "test")
    cmd.Dir = filepath.Dir(cm.src)
    output, err := cmd.CombinedOutput()
    if err == nil {
        cm.failed++
        cm.logf("mutation %d failed to break any tests", n)
    } else if _, ok := err.(*exec.ExitError); ok {
        lines := bytes.Split(output, []byte("\n"))
        lastLine := lines[len(lines)-2]
        if bytes.HasPrefix(lastLine, []byte("FAIL")) {
            cm.logf("mutation %d broke the tests properly", n)
        } else {
            cm.logf("mutation %d created an error: %s", n, lastLine)
        }
    } else {
        return fmt.Errorf("mutation %d failed to run: %s", n, err)
    }
    return nil
}
