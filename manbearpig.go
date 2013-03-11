package main

import (
    "flag"
    "fmt"
    "github.com/darkhelmet/manbearpig/mutations"
    "log"
    "os"
    "path/filepath"
)

var (
    list   = flag.Bool("list", false, "Show available mutations")
    logger = log.New(os.Stdout, "", log.LstdFlags)
)

func main() {
    flag.Parse()
    if *list {
        for name, _ := range mutations.Mutations {
            fmt.Println(name)
        }
        os.Exit(0)
    }

    m := EnsureValidMutation()
    pkg := EnsureValidPackage()
    tmp := EnsureTmpDir()

    logger.Printf("mutating in %s", tmp)

    CopyFiles(pkg.Dir, tmp)

    for _, f := range pkg.GoFiles {
        src := filepath.Join(tmp, f)
        m.Prepare(src, logger)
        m.Run()
    }

    os.Exit(m.FailureCount())
}
