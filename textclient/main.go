package main

import (
	"fmt"
	"os"
	"runtime/debug"
)

func main() {
	var client Client
	initialized := false

	defer func() {
		if initialized {
			client.Deinit() // resets terminal changes
		}

		if r := recover(); r == nil {
			fmt.Println("Thanks for playing :)")
			os.Exit(0)
		} else {
			_, _ = fmt.Fprintln(os.Stderr, r)
			debug.PrintStack()
			os.Exit(2)
		}
	}()

	if err := client.Init(); err == nil {
		initialized = true
	} else {
		_, _ = fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}

	if err := client.Run(); err != nil {
		client.Deinit() // resets terminal changes
		_, _ = fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
}
