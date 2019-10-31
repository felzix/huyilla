package main

import (
	"github.com/op/go-logging"
	"os"
)

func main() {
	logging.SetBackend(
		logging.AddModuleLevel(
			logging.NewLogBackend(os.Stderr, "", 0))).SetLevel(logging.DEBUG, "")

	guiClient := NewGuiClient()
	addLight(guiClient.rootScene)
	addAxes(guiClient.rootScene)
	setBackgroundColor(guiClient.app)

	if err := guiClient.Run(); err != nil {
		panic(err)
	}
}
