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

	var webServerError chan error
	go engine.Serve(webServerError)
mainloop:
	for {
		select {
		case <-sigs:
			break mainloop
		case err := <-webServerError:
			fail(err)
		case <-time.After(time.Millisecond * 500):
			if err := engine.Tick(); err != nil {
				fail(err)
			}
		}
	}

	fmt.Println("Engine stopped!")
}

func fail(err error) {
	os.Stderr.WriteString("Error!\n")
	os.Stderr.WriteString(err.Error())
	os.Exit(1)
}
