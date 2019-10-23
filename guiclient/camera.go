package main

import (
	"fmt"
	"github.com/g3n/engine/camera"
	"github.com/g3n/engine/core"
	"github.com/g3n/engine/gui"
	"github.com/g3n/engine/math32"
	"github.com/g3n/engine/window"
)

type CameraController struct {
	core.Dispatcher

	cam *camera.Camera

	RotSpeed float32
}

func NewCameraController(cam *camera.Camera) {
	controller := new(CameraController)
	controller.Dispatcher.Initialize()
	controller.cam = cam
	controller.RotSpeed = 3.0
	gui.Manager().SetCursorFocus(controller)
	controller.SubscribeID(window.OnCursor, &controller, controller.onCursor)
	// gui.Manager().SubscribeID(window.OnCursor, &controller, controller.onCursor)
}

func (controller *CameraController) Dispose() {
	controller.UnsubscribeID(window.OnCursor, &controller)
}

// winSize returns the window height or width based on the camera reference axis.
func (controller *CameraController) winSize() float32 {
	width, size := window.Get().GetSize()
	if controller.cam.Axis() == camera.Horizontal {
		size = width
	}
	return float32(size)
}

var cursorCalls = 0

func (controller *CameraController) onCursor(eventName string, event interface{}) {
	cursorEvent := event.(*window.CursorEvent)

	width, height := window.Get().GetSize()
	c := -2 * math32.Pi * controller.RotSpeed
	cX := c / float32(width)
	cY := c / float32(height)
	x0, y0 := MiddleOfScreen()
	xDelta := cursorEvent.Xpos - float32(x0)
	yDelta := cursorEvent.Ypos - float32(y0)


	if math32.Abs(xDelta) >= 1 || math32.Abs(yDelta) >= 1 {
		cursorCalls++
		fmt.Println("cursor calls:", cursorCalls)
		// fmt.Println("window", width, height)
		// fmt.Println("middle", x0, y0)
		// fmt.Println("cursor", cursorEvent.Xpos, cursorEvent.Ypos)
		// fmt.Println("delta", xDelta, yDelta)

		controller.Rotate(cX*xDelta, cY*yDelta)
		SetCursorPos(x0, y0)
	}
}

// Rotate rotates the camera in place.
func (controller *CameraController) Rotate(xDelta, yDelta float32) {
	p := controller.cam.Rotation()
	// mixing x and y is intended here
	x := p.X + yDelta
	y := p.Y + xDelta
	controller.cam.SetRotation(x, y, 0)
}

func SetCursorPos(x, y float64) {
	w := window.Get()
	// TODO use a type switch or fork the lib to expand the interface
	//      See https://w3c.github.io/pointerlock if choosing to fork
	gw := w.(*window.GlfwWindow)
	gw.SetCursorPos(x, y - 1) // sets it to y+1 for some reason, so subtract one
}

func MiddleOfScreen() (float64, float64) {
	width, height := window.Get().GetSize()
	var halfWidth, halfHeight float64

	halfWidth = float64(width) / 2
	halfHeight = float64(height) / 2

	return halfWidth, halfHeight
}
