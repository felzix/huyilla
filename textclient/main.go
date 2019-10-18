package main

import (
	"fmt"
	"os"
	"runtime/debug"
)

func main() {
	backend, closeFn, err := MakeFileLogBackend("/tmp/huyilla-log")
	if err != nil {
		panic(err)
	}
	defer func() {
		if err := closeFn(); err != nil {
			panic(err)
		}
	}()
	Log.SetBackend(backend)

	var client TextClient
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
		initialized = false
		_, _ = fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
}
