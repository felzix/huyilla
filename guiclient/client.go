package main

import (
	"fmt"
	"github.com/felzix/huyilla/client"
	C "github.com/felzix/huyilla/constants"
	"github.com/felzix/huyilla/types"
	g3nApp "github.com/g3n/engine/app"
	g3nCamera "github.com/g3n/engine/camera"
	"github.com/g3n/engine/core"
	"github.com/g3n/engine/gls"
	"github.com/g3n/engine/graphic"
	"github.com/g3n/engine/gui"
	"github.com/g3n/engine/light"
	"github.com/g3n/engine/math32"
	"github.com/g3n/engine/renderer"
	"github.com/g3n/engine/util/helper"
	"github.com/g3n/engine/window"
	"github.com/gdamore/tcell"
	"sync"
	"time"
)

type GuiClient struct {
	sync.Mutex

	Cache    *Cache
	player   *types.PlayerDetails
	username string
	api      *client.API

	displayRadius uint

	quitq    chan struct{}
	quitOnce sync.Once
	err      error

	eventq chan tcell.Event

	app        *g3nApp.Application
	rootScene  *core.Node
	camera     *g3nCamera.Camera
	playerNode *graphic.Mesh
}

func NewGuiClient() *GuiClient {
	cam := g3nCamera.NewPerspective(1, 1, 1000, 45, 0)
	app, scene := setupGraphics(cam)

	NewCameraController(cam)

	guiClient := &GuiClient{
		Cache:    NewCache(scene),
		player:   nil,
		username: "felzix",
		api:      nil,

		displayRadius: 3,

		quitq:    make(chan struct{}),
		quitOnce: sync.Once{},
		err:      nil,

		eventq: make(chan tcell.Event),

		app:       app,
		rootScene: scene,
		camera:    cam,
	}

	app.Subscribe(window.OnKeyUp, func(_ string, event interface{}) {
		keyEvent := event.(*window.KeyEvent)
		switch keyEvent.Key {
		case window.KeyEscape:
			guiClient.app.Exit()
		}
	})

	app.Subscribe(window.OnChar, func(_ string, event interface{}) {
		charEvent := event.(*window.CharEvent)

		switch charEvent.Char {
		// TODO base move commands on player's rotation

		case 'w':
			target := guiClient.player.Entity.Location.Derive(1, 0, 0, C.CHUNK_SIZE)
			fmt.Println(guiClient.player.Entity.Location.Voxel.ToString(), "->", target.ToString())
			if err := guiClient.api.IssueMoveAction(target); err != nil {
				panic(err)
			}
		case 's':
			target := guiClient.player.Entity.Location.Derive(-1, 0, 0, C.CHUNK_SIZE)
			fmt.Println(guiClient.player.Entity.Location.Voxel.ToString(), "->", target.ToString())
			if err := guiClient.api.IssueMoveAction(target); err != nil {
				panic(err)
			}
		case 'a':
			target := guiClient.player.Entity.Location.Derive(0, 1, 0, C.CHUNK_SIZE)
			fmt.Println(guiClient.player.Entity.Location.Voxel.ToString(), "->", target.ToString())
			if err := guiClient.api.IssueMoveAction(target); err != nil {
				panic(err)
			}
		case 'd':
			target := guiClient.player.Entity.Location.Derive(0, -1, 0, C.CHUNK_SIZE)
			fmt.Println(guiClient.player.Entity.Location.Voxel.ToString(), "->", target.ToString())
			if err := guiClient.api.IssueMoveAction(target); err != nil {
				panic(err)
			}
		}
	})

	return guiClient
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
		case <-time.After(time.Millisecond * 500): // poll engine only so often
			if guiClient.api == nil || guiClient.player == nil {
				continue // user is still entering in their information
			}

			// ignores error because getting world age is eqiuvalent to querying the readiness of the server
			age, _ := guiClient.api.GetWorldAge()

			if age > guiClient.Cache.GetAge() {
				guiClient.Cache.SetAge(age)

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

				guiClient.Cache.SetChunks(chunks)
				guiClient.Cache.Draw(guiClient)
			}
		}
	}
}

func setupGraphics(cam *g3nCamera.Camera) (*g3nApp.Application, *core.Node) {
	// Create application and scene
	app := g3nApp.App()
	scene := core.NewNode()

	// Set the scene to be managed by the gui manager
	gui.Manager().Set(scene)

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

	return app, scene
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
