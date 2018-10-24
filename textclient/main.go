package main

import (
    "fmt"
    "os"
)


func main () {
    var client Client
    if err := client.Init(); err != nil {
        fmt.Fprintf(os.Stderr, "%s\n", err)
        os.Exit(1)
    }
    defer client.Deinit()  // resets terminal changes

    if err := client.Run(); err != nil {
        fmt.Fprintf(os.Stderr, "%s\n", err)
        os.Exit(1)
    }
}
