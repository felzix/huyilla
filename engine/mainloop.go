package main

import (
	"fmt"
	"github.com/felzix/huyilla/engine/engine"
	"os"
	"os/signal"
	"syscall"
	"time"
)



const PORT = 8080

func main() {
	engine.Log.Info("Starting engine...")

	huyilla := &engine.Engine{}
	if err := huyilla.Init("/tmp/huyilla"); err != nil {
		panic(err)
	}

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	engine.Log.Infof("Engine started @ port %d!\n", PORT)

	var webServerError chan error
	huyilla.Serve(PORT, webServerError)

mainloop:
	for {
		select {
		case <-sigs:
			break mainloop
		case err := <-webServerError:
			fail(err)
		case <-time.After(time.Millisecond * 50): // tick period
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
