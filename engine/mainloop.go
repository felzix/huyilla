package main

import (
	"fmt"
	"github.com/felzix/huyilla/engine/engine"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	fmt.Println("Starting engine...")

	huyilla := &engine.Engine{}
	if err := huyilla.Init("/tmp/huyilla"); err != nil {

	}

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	fmt.Println("Engine started!")

	var webServerError chan error
	huyilla.Serve(8080, webServerError)

mainloop:
	for {
		select {
		case <-sigs:
			break mainloop
		case err := <-webServerError:
			fail(err)
		case <-time.After(time.Millisecond * 500):
			if err := huyilla.Tick(); err != nil {
				fail(err)
			}
		}
	}

	fmt.Println("Engine stopped!")
}

func fail(err error) {
	_, _ = os.Stderr.WriteString("Error!\n")
	_, _ = os.Stderr.WriteString(err.Error())
	os.Exit(1)
}
