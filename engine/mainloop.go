package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	fmt.Println("Engine started!")

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	for {
		select {
		case <-sigs:
			break
		case <-time.After(time.Millisecond * 50):
			// TODO tick
		}
	}

	fmt.Println("Engine stopped!")
}
