package main

import (
	"github.com/felzix/huyilla/client"
	C "github.com/felzix/huyilla/constants"
	"github.com/felzix/huyilla/content"
	"github.com/felzix/huyilla/types"
	g3nApp "github.com/g3n/engine/app"
	g3nCamera "github.com/g3n/engine/camera"
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
	"github.com/gdamore/tcell"
	"sync"
	"time"
)

func buildVoxels(scene *core.Node, chunk *types.DetailedChunk, offset *types.Point) {
	for x := 0; x < C.CHUNK_SIZE; x++ {
		for y := 0; y < C.CHUNK_SIZE; y++ {
			for z := 0; z < C.CHUNK_SIZE; z++ {
				voxel := chunk.GetVoxel(uint64(x), uint64(y), uint64(z))
				if isDrawn(voxel) {
					trueX := float32(x + int(offset.X * 16))
					trueY := float32(y + int(offset.Y * 16))
					trueZ := float32(z + int(offset.Z * 16))
					makeVoxel(scene, trueX, trueY, trueZ, voxel)
				}
			}
		}
	}
	for _, e := range chunk.Entities {
		eX := float32(e.Location.Voxel.X + (offset.X * 16))
		eY := float32(e.Location.Voxel.Y + (offset.Y * 16))
		eZ := float32(e.Location.Voxel.Z + (offset.Z * 16))
		makeEntity(scene, eX, eY, eZ, e)
	}
}

func makeEntity(scene *core.Node, x, y, z float32, entity *types.Entity) {
	def := content.EntityDefinitions[entity.Type]
	geom := geometries[def.Form]
	mat := materials[def.Material]

	mesh := graphic.NewMesh(geom, mat)
	mesh.SetPosition(x, y, z + 1)
	scene.Add(mesh)
	mesh.SetRotation(math32.Pi/2, 0, 0)
}

func isDrawn(voxel types.Voxel) bool {
	M := content.MATERIAL
	v := voxel.Expand()
	switch v.Material {
	case M["air"]:
		return false
	case M["dirt"]:
		return true
	case M["grass"]:
		return true
	case M["water"]:
		return true
	default:
		return false
	}
}

var geometries = make(map[uint64]geometry.IGeometry)

func makeGeometries() {
	geometries[content.FORM["cube"]] = geometry.NewCube(1)
	geometries[content.FORM["human"]] = geometry.NewCylinder(0.3, 1.8, 16, 1, true, true)
}

var materials = make(map[uint64]material.IMaterial)

func makeMaterials() {
	materials[content.MATERIAL["dirt"]] = material.NewStandard(math32.NewColor("SaddleBrown"))
	materials[content.MATERIAL["grass"]] = material.NewStandard(math32.NewColor("SpringGreen"))
	materials[content.MATERIAL["water"]] = material.NewStandard(math32.NewColor("DarkBlue"))

	materials[content.MATERIAL["human"]] = material.NewStandard(math32.NewColor("DarkRed"))
}

func makeVoxel(scene *core.Node, x, y, z float32, voxel types.Voxel) {
	v := voxel.Expand()
	geom := geometries[v.Form]
	mat := materials[v.Material]

	mesh := graphic.NewMesh(geom, mat)
	scene.Add(mesh)
	mesh.SetPosition(x, y, z)
}

type GuiClient struct {
	sync.Mutex

	world    *client.WorldCache
	player   *types.PlayerDetails
	username string
	api      *client.API

	displayRadius  uint

	quitq    chan struct{}
	quitOnce sync.Once
	err      error

	eventq chan tcell.Event

	app *g3nApp.Application
	rootScene *core.Node
	camera *g3nCamera.Camera
}

func NewGuiClient() *GuiClient {
	app, scene, cam := setupGraphics()

	return &GuiClient{
		world: client.NewWorldCache(),
		player: nil,
		username: "felzix",
		api: nil,

		displayRadius: 3,

		quitq: make(chan struct{}),
		quitOnce: sync.Once{},
		err: nil,

		eventq:  make(chan tcell.Event),

		app: app,
		rootScene: scene,
		camera: cam,
	}
}

func (guiClient *GuiClient) Run() error {
	if err := guiClient.Auth(); err != nil {
		return err
	}
	go guiClient.EnginePoller()
	guiClient.runGraphics()

	return guiClient.err
}

func (guiClient *GuiClient) Auth() error {
	guiClient.api = client.NewAPI("http://localhost:8080", guiClient.username, "murakami")

	if exists, err := guiClient.api.UserExists(); err == nil {
		if !exists {
			if err := guiClient.api.Signup(); err != nil {
				return err
			}
		}
	} else {
		return err
	}

	if err := guiClient.api.Login(); err != nil {
		return err
	}

	entity, err := guiClient.api.GetPlayer(guiClient.username)
	if err != nil {
		return err
	}

	guiClient.player = &types.PlayerDetails{
		Player: &types.Player{
			Name:     guiClient.username,
			EntityId: entity.Id,
		},
		Entity: entity,
	}

	return nil
}

func (guiClient *GuiClient) runGraphics() {
	guiClient.app.Run(func(renderer *renderer.Renderer, deltaTime time.Duration) {
		guiClient.app.Gls().Clear(gls.DEPTH_BUFFER_BIT | gls.STENCIL_BUFFER_BIT | gls.COLOR_BUFFER_BIT)
		if err := renderer.Render(guiClient.rootScene, guiClient.camera); err != nil {
			panic(err)
		}
	})
}


func (guiClient *GuiClient) Quit(err error) {
	guiClient.err = err
	guiClient.quitOnce.Do(func() {
		close(guiClient.quitq)
	})
}

func (guiClient *GuiClient) EnginePoller() {
	for {
		select {
		case <-guiClient.quitq:
			return
		case <-time.After(time.Millisecond * 2000): // poll engine only so often
			if guiClient.api == nil || guiClient.player == nil {
				continue // user is still entering in their information
			}

			// ignores error because getting world age is eqiuvalent to querying the readiness of the server
			age, _ := guiClient.api.GetWorldAge()

			if age > guiClient.world.GetAge() {
				guiClient.world.SetAge(age)

				entity, err := guiClient.api.GetPlayer(guiClient.username)
				if err != nil {
					guiClient.Quit(err)
					return
				}

				guiClient.player.Entity = entity

				center := guiClient.player.Entity.Location.Chunk

				chunks, err := guiClient.api.GetChunks(center, C.ACTIVE_CHUNK_RADIUS)
				if err != nil {
					guiClient.Quit(err)
					return
				}

				for i, chunk := range chunks.Chunks {
					point := chunks.Points[i]
					guiClient.world.SetChunk(point, chunk)
				}
				buildVoxels(guiClient.rootScene, guiClient.world.GetChunk(center), &types.Point{})
				below := types.NewPoint(center.X, center.Y, center.Z - 1)
				buildVoxels(guiClient.rootScene, guiClient.world.GetChunk(below), &types.Point{Z: -1})
			}
		}
	}
}

func setupGraphics() (*g3nApp.Application, *core.Node, *g3nCamera.Camera) {
	// Create application and scene
	app := g3nApp.App()
	scene := core.NewNode()

	// Set the scene to be managed by the gui manager
	gui.Manager().Set(scene)

	cam := addCamera(scene)

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

	makeGeometries()
	makeMaterials()

	return app, scene, cam
}

func addCamera(scene *core.Node) *g3nCamera.Camera {
	// Create perspective camera
	cam := g3nCamera.NewPerspective(1, 1, 1000, 45, 0)
	// cam := g3nCamera.New(1)
	cam.SetPosition(0, 0, 1)
	scene.Add(cam)
	cam.SetRotation(1.5, 0, 0)

	// Set up orbit control for the camera
	// control := g3nCamera.NewOrbitControl(cam)

	return cam
}

// Create and add lights to the scene
func addLight(scene *core.Node) {
	scene.Add(light.NewAmbient(&math32.Color{R: 1.0, G: 1.0, B: 1.0}, 0.8))
	pointLight := light.NewPoint(&math32.Color{R: 1, G: 1, B: 1}, 5.0)
	pointLight.SetPosition(1, 0, 2)
	scene.Add(pointLight)
}

// Create and add an axis helper to the scene
func addAxes(scene *core.Node) {
	scene.Add(helper.NewAxes(0.5))
}

// Set background color to gray
func setBackgroundColor(app *g3nApp.Application) {
	app.Gls().ClearColor(0.5, 0.5, 0.5, 1.0)
}

func main() {
	guiClient := NewGuiClient()
	addLight(guiClient.rootScene)
	addAxes(guiClient.rootScene)
	setBackgroundColor(guiClient.app)

	if err := guiClient.Run(); err != nil {
		panic(err)
	}
}
