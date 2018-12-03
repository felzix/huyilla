package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	fmt.Println("Starting engine...")

	engine := &Engine{}
	engine.Init()

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	fmt.Println("Engine started!")

mainloop:
	for {
		select {
		case <-sigs:
			break mainloop
		case <-time.After(time.Millisecond * 500):
			if err := engine.Tick(); err != nil {
				os.Stderr.WriteString("Error!\n")
				os.Stderr.WriteString(err.Error())
				os.Exit(1)
			}
		}
	}

	fmt.Println("Engine stopped!")
}
