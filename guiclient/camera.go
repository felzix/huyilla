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
}

func NewCameraController(cam *camera.Camera) {
	controller := new(CameraController)
	controller.Dispatcher.Initialize()
	controller.cam = cam
	controller.RotSpeed = 2.0
	gui.Manager().SubscribeID(window.OnCursor, &controller, controller.onCursor)
	// controller.SubscribeID(window.OnCursor, &controller, controller.onCursor)
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

func (controller *CameraController) onCursor (eventName string, event interface{}) {
	w := window.Get()
	// TODO use a type switch or fork the lib to expand the interface
	//      See https://w3c.github.io/pointerlock if choosing to fork
	gw := w.(*window.GlfwWindow)
	width, height := w.GetSize()

	cursorEvent := event.(*window.CursorEvent)
	cX := -2 * math32.Pi * controller.RotSpeed / float32(width)
	cY := -2 * math32.Pi * controller.RotSpeed / float32(height)
	x0, y0 := float64(width)/2, float64(height)/2
	xDelta := cursorEvent.Xpos - float32(x0)
	yDelta := cursorEvent.Ypos - float32(y0)
	controller.Rotate(cX * xDelta, cY * yDelta)
	gw.SetCursorPos(x0, y0)
}

// Rotate rotates the camera in place.
func (controller *CameraController) Rotate(xDelta, yDelta float32) {
	p := controller.cam.Rotation()
	// mixing x and y is intended here
	x := p.X + yDelta
	y := p.Y + xDelta
	controller.cam.SetRotation(x, y, 0)
}
