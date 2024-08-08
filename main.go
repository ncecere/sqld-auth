package main

import (
    "fmt"
    "os"

    "github.com/ncecere/sqld-auth/cmd"
)

func main() {
    if err := cmd.Execute(); err != nil {
        fmt.Println(err)
        os.Exit(1)
    }
}