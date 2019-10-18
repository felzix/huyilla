package main

func main() {
	guiClient := NewGuiClient()
	addLight(guiClient.rootScene)
	addAxes(guiClient.rootScene)
	setBackgroundColor(guiClient.app)

	if err := guiClient.Run(); err != nil {
		panic(err)
	}
}
