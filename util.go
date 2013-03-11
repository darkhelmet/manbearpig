package main

import (
    "flag"
    "github.com/darkhelmet/manbearpig/mutations"
    "go/build"
    "io"
    "io/ioutil"
    "os"
    "path/filepath"
)

var (
    mutation   = flag.String("mutation", "", "The mutation to perform")
    importPath = flag.String("import", "", "The import path to mutate")
)

func EnsureValidMutation() mutations.Mutation {
    if *mutation == "" {
        logger.Fatalf("no mutation specified")
    }

    m, ok := mutations.Mutations[*mutation]
    if !ok {
        logger.Fatalf("%#v is not a valid mutation", *mutation)
    }
    return m
}

func EnsureValidPackage() *build.Package {
    if *importPath == "" {
        logger.Fatalf("no import path specified")
    }

    pkg, err := build.Import(*importPath, "", 0)
    if err != nil {
        logger.Fatalf("failed to import package: %s", err)
    }
    return pkg
}

func EnsureTmpDir() string {
    tmp, err := ioutil.TempDir("", "manbearpig")
    if err != nil {
        logger.Fatalf("failed to create tmp directory: %s", err)
    }
    return tmp
}

func CopyFile(src, dir string) error {
    name := filepath.Base(src)
    srcFile, err := os.Open(src)
    if err != nil {
        return err
    }
    defer srcFile.Close()

    dstFile, err := os.Create(filepath.Join(dir, name))
    if err != nil {
        return err
    }
    defer dstFile.Close()

    _, err = io.Copy(dstFile, srcFile)
    return err
}

func CopyFiles(src, dst string) {
    contents, err := ioutil.ReadDir(src)
    if err != nil {
        logger.Fatalf("failed reading directory: %s", err)
    }
    for _, f := range contents {
        if f.Mode()&os.ModeType == 0 {
            err := CopyFile(filepath.Join(src, f.Name()), dst)
            if err != nil {
                logger.Fatalf("failed copying %s: %s", f.Name(), err)
            }
        }
    }
}
