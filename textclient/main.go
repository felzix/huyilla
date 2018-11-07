package main

import (
    "fmt"
    "os"
)


func main () {
    var client Client

    defer func() {
        if r := recover(); r != nil {
            fmt.Fprintf(os.Stderr, "Recovered %v", r)
            os.Exit(2)
        }
    }()

    defer finish(&client, 0, nil)

    if err := client.Init(); err != nil {
        finish(&client, 1, err)
    }

    if err := client.Run(); err != nil {
        finish(&client, 1, err)
    }

    defer fmt.Println("Thanks for playing!")
}

func finish(client *Client, returnCode int, err error) {
    defer os.Exit(returnCode)

    defer client.Deinit()  // resets terminal changes

    if err == nil {
        fmt.Println("Thanks for playing!")
    } else {
        fmt.Fprintf(os.Stderr, "%s\n", err)
    }
}
