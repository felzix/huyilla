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
	CursorPosition math32.Vector2
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
	w2 := w.(*window.GlfwWindow)
	width, height := w.GetSize()

	cursorEvent := event.(*window.CursorEvent)
	cX := -2 * math32.Pi * controller.RotSpeed / float32(width)
	cY := -2 * math32.Pi * controller.RotSpeed / float32(height)
	xDelta := cursorEvent.Xpos - controller.CursorPosition.X
	yDelta := cursorEvent.Ypos - controller.CursorPosition.Y
	controller.Rotate(cX * xDelta, cY * yDelta)

	controller.CursorPosition.Set(float32(width)/2, float32(height)/2)
	w2.SetCursorPos(float64(width)/2, float64(height)/2)
}

// Rotate rotates the camera in place.
func (controller *CameraController) Rotate(xDelta, yDelta float32) {
	p := controller.cam.Rotation()
	// mixing x and y is intended here
	x := p.X + yDelta
	y := p.Y + xDelta
	controller.cam.SetRotation(x, y, 0)
}
