package main

import (
	"fmt"
	"github.com/akamensky/argparse"
	"github.com/felzix/huyilla/constants"
	"github.com/felzix/huyilla/engine/engine"
	uuid "github.com/satori/go.uuid"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type Options struct {
	DBType string
	Port   int
}

func parse() (*Options, error) {
	parser := argparse.NewParser("Huyilla Engine", "")
	dbType := parser.Selector("D", "database", []string{"diskv", "memory"}, &argparse.Options{
		Help:    "use diskv for persistence",
		Default: "diskv",
	})
	port := parser.Int("p", "port", &argparse.Options{
		Default: 8080,
	})
	err := parser.Parse(os.Args)
	if err != nil {
		fmt.Print(parser.Usage(err))
		return nil, err
	}

	return &Options{
		DBType: *dbType,
		Port:   *port,
	}, nil
}

func main() {
	opts, err := parse()
	if err != nil {
		return // only errors on bad usage and help was already printed
	}

	var db engine.Database
	switch opts.DBType {
	case "diskv":
		db = engine.NewDisKVDatabase(
			"/tmp/huyilla",
			"/tmp/tempdir-huyilla-"+uuid.NewV4().String(),
			16*1024*1024) // 16 MB
	case "memory":
		db = engine.NewMemoryDatabase()
	}

	engine.Log.Info("Starting engine...")

	huyilla, err := engine.NewEngine(
		constants.SEED,
		engine.NewGrassWorldGenerator(),
		//engine.NewLakeWorldGenerator(3),
		db)

	if err != nil {
		panic(err)
	}

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	engine.Log.Infof("Engine started @ port %d!\n", opts.Port)

	var webServerError chan error
	huyilla.Serve(opts.Port, webServerError)

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
