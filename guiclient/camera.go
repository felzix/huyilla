package main

import (
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
	LastPosition math32.Vector2
}

func NewCameraController(cam *camera.Camera) {
	controller := new(CameraController)
	controller.Dispatcher.Initialize()
	controller.cam = cam
	controller.LastPosition = math32.Vector2{}
	controller.RotSpeed = 3.0
	gui.Manager().SetCursorFocus(controller)
	controller.SubscribeID(window.OnCursor, &controller, controller.onCursor)
	SetCursorInputMode(window.CursorDisabled)
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

func (controller *CameraController) onCursor(eventName string, event interface{}) {
	cursorEvent := event.(*window.CursorEvent)

	width, height := window.Get().GetSize()
	c := -2 * math32.Pi * controller.RotSpeed
	cX := c / float32(width)
	cY := c / float32(height)
	x0 := controller.LastPosition.X
	y0 := controller.LastPosition.Y
	controller.LastPosition.X = cursorEvent.Xpos
	controller.LastPosition.Y = cursorEvent.Ypos

	if x0 == 0 && y0 == 0 {
		return // don't rotate camera on first event
	}

	xDelta := cursorEvent.Xpos - x0
	yDelta := cursorEvent.Ypos - y0

	controller.Rotate(cX*xDelta, cY*yDelta)
}

// Rotate rotates the camera in place.
func (controller *CameraController) Rotate(xDelta, yDelta float32) {
	p := controller.cam.Rotation()
	// mixing x and y is intended here
	x := p.X + yDelta
	y := p.Y + xDelta
	controller.cam.SetRotation(x, y, 0)
}

func SetCursorInputMode(mode window.CursorMode) {
	w := window.Get()
	// TODO use a type switch or fork the lib to expand the interface
	//      See https://w3c.github.io/pointerlock if choosing to fork
	gw := w.(*window.GlfwWindow)
	gw.SetInputMode(window.CursorInputMode, mode)
}
