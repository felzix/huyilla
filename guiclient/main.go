package main

import (
	C "github.com/felzix/huyilla/constants"
	"github.com/felzix/huyilla/content"
	"github.com/felzix/huyilla/engine/engine"
	"github.com/felzix/huyilla/types"
	g3nApp "github.com/g3n/engine/app"
	"github.com/g3n/engine/camera"
	"github.com/g3n/engine/core"
	"github.com/g3n/engine/geometry"
	"github.com/g3n/engine/gls"
	"github.com/g3n/engine/graphic"
	"github.com/g3n/engine/gui"
	"github.com/g3n/engine/light"
	"github.com/g3n/engine/material"
	"github.com/g3n/engine/math32"
	"github.com/g3n/engine/renderer"
	"github.com/g3n/engine/util/helper"
	"github.com/g3n/engine/window"
	"time"
)

func makeChunk() *types.Chunk {
	world := &engine.World{Seed: C.SEED}

	if err := world.Init("/tmp/huyilla-gui", 16*1024*1024); err != nil { // 16 MB
		panic(err)
	}

	center := types.NewPoint(0, 0, 0)
	chunk, err := world.Chunk(center)
	if err != nil {
		panic(err)
	}

	return chunk
}

func voxelToColor(voxel types.Voxel) (string, bool) {
	M := content.MATERIAL

	v := voxel.Expand()

	// see /Users/robertdavidson/go/src/github.com/g3n/engine/math32/color.go
	switch v.Material {
	case M["air"]:
		return "", false
	case M["dirt"]:
		return "SaddleBrown", true
	case M["grass"]:
		return "SpringGreen", true
	case M["water"]:
		return "DarkBlue", true
	default:
		return "", false
	}
}

func buildVoxels(scene *core.Node, chunk *types.Chunk) {
	for x := 0; x < C.CHUNK_SIZE; x++ {
		for y := 0; y < C.CHUNK_SIZE; y++ {
			for z := 0; z < C.CHUNK_SIZE; z++ {
				voxel := chunk.GetVoxel(uint64(x), uint64(y), uint64(z))
				color, drawn := voxelToColor(voxel)
				if drawn {
					makeVoxel(scene, float32(x), float32(y), float32(z), color)
				}
			}
		}
	}
}

func makeVoxel(scene *core.Node, x, y, z float32, color string) {
	geom := geometry.NewCube(1)
	mat := material.NewStandard(math32.NewColor(color))
	mesh := graphic.NewMesh(geom, mat)
	mesh.SetPosition(x, y, z)
	scene.Add(mesh)
}

func main() {
	// Create application and scene
	app := g3nApp.App()
	scene := core.NewNode()

	// Set the scene to be managed by the gui manager
	gui.Manager().Set(scene)

	// Create perspective camera
	cam := camera.New(1)
	cam.SetPosition(0, 0, 4)
	scene.Add(cam)

	// Set up orbit control for the camera
	camera.NewOrbitControl(cam)

	// Set up callback to update viewport and camera aspect ratio when the window is resized
	onResize := func(evname string, ev interface{}) {
		// Get framebuffer size and update viewport accordingly
		width, height := app.GetSize()
		app.Gls().Viewport(0, 0, int32(width), int32(height))
		// Update the camera's aspect ratio
		cam.SetAspect(float32(width) / float32(height))
	}
	app.Subscribe(window.OnWindowSize, onResize)
	onResize("", nil)

	buildVoxels(scene, makeChunk())

	// Create and add app button to the scene
	// btn := gui.NewButton("Make Red")
	// btn.SetPosition(100, 40)
	// btn.SetSize(40, 40)
	// btn.Subscribe(gui.OnClick, func(name string, ev interface{}) {
	// 	mat.SetColor(math32.NewColor("DarkRed"))
	// })
	// scene.Add(btn)

	// Create and add lights to the scene
	scene.Add(light.NewAmbient(&math32.Color{1.0, 1.0, 1.0}, 0.8))
	pointLight := light.NewPoint(&math32.Color{1, 1, 1}, 5.0)
	pointLight.SetPosition(1, 0, 2)
	scene.Add(pointLight)

	// Create and add an axis helper to the scene
	scene.Add(helper.NewAxes(0.5))

	// Set background color to gray
	app.Gls().ClearColor(0.5, 0.5, 0.5, 1.0)

	// Run the application
	app.Run(func(renderer *renderer.Renderer, deltaTime time.Duration) {
		app.Gls().Clear(gls.DEPTH_BUFFER_BIT | gls.STENCIL_BUFFER_BIT | gls.COLOR_BUFFER_BIT)
		if err := renderer.Render(scene, cam); err != nil {
			panic(err)
		}
	})
}
