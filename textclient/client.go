package main

import (
	"github.com/felzix/huyilla/client"
	"github.com/felzix/huyilla/constants"
	"sync"
	"time"

	react "github.com/felzix/go-curses-react"
	"github.com/felzix/huyilla/types"
	"github.com/gdamore/tcell"
)

type TextClient struct {
	sync.Mutex

	world    *client.WorldCache
	player   *types.PlayerDetails
	username string
	api      *client.API

	screen         *react.Screen
	viewDepthDelta int
	displayRadius  uint

	quitq    chan struct{}
	quitOnce sync.Once
	err      error

	eventq chan tcell.Event
}

func (textClient *TextClient) Init() error {
	textClient.displayRadius = 15

	if screen, err := react.NewScreen(); err == nil {
		textClient.screen = screen
	} else {
		return err
	}

	textClient.quitq = make(chan struct{})
	textClient.eventq = make(chan tcell.Event)

	root := MakeApp()
	root.Props = react.Properties{
		"textClient": textClient,
	}
	if err := textClient.screen.Init(root, func(err error) error {
		textClient.Quit(err)
		return nil
	}); err != nil {
		return err
	}

	textClient.world = client.NewWorldCache()

	return nil
}

func (textClient *TextClient) Deinit() {
	textClient.screen.TCellScreen.Fini()
}

func (textClient *TextClient) Run() error {
	go textClient.EventPoller()  // gets user input
	go textClient.EnginePoller() // gets world state from the engine

loop:
	for {
		if err := textClient.Draw(); err != nil {
			textClient.Quit(err)
			break loop
		}
		select {
		case <-textClient.quitq:
			break loop
		// draw-input loop runs no faster than once every 10ms
		case <-time.After(time.Millisecond * 10):
		case ev := <-textClient.eventq:
			if err := textClient.HandleEvent(ev); err != nil {
				textClient.Quit(err)
				break loop
			}
		}
	}

	// Inject a wakeup interrupt
	iev := tcell.NewEventInterrupt(nil)
	if err := textClient.screen.TCellScreen.PostEvent(iev); err != nil {
		return err
	}

	return textClient.err
}

func (textClient *TextClient) HandleEvent(e tcell.Event) error {
	switch e := e.(type) {
	case *tcell.EventResize:
		textClient.handleResize(e)
	case *tcell.EventKey:
		return textClient.handleKey(e)
	case *tcell.EventMouse:
		return textClient.handleMouse(e)
	}

	return nil
}

func (textClient *TextClient) handleResize(e *tcell.EventResize) {
	textClient.screen.Resize()
	textClient.screen.TCellScreen.Sync() // visually jarring but needed after a resize
}

func (textClient *TextClient) handleKey(e *tcell.EventKey) error {
	textClient.Lock()
	defer textClient.Unlock()
	return textClient.screen.HandleKey(e)
}

func (textClient *TextClient) handleMouse(e *tcell.EventMouse) error {
	// x, y := e.Position()
	// fmt.Printf("(%d,%d)", x, y)
	return nil
}

func (textClient *TextClient) Draw() error {
	textClient.Lock()
	defer textClient.Unlock()

	err := textClient.screen.Draw()
	if err == nil {
		textClient.screen.TCellScreen.Show()
	}
	return err
}

func (textClient *TextClient) EventPoller() {
	for {
		select {
		case <-textClient.quitq:
			return
		default:
		}

		e := textClient.screen.TCellScreen.PollEvent()
		if e == nil {
			return
		}

		select {
		case <-textClient.quitq:
			return
		case textClient.eventq <- e:
		}
	}
}

func (textClient *TextClient) EnginePoller() {
	for {
		select {
		case <-textClient.quitq:
			return
		case <-time.After(time.Millisecond * 50): // poll engine only so often
			if textClient.api == nil || textClient.player == nil {
				continue // user is still entering in their information
			}

			// ignores error because getting world age is eqiuvalent to querying the readiness of the server
			age, _ := textClient.api.GetWorldAge()

			if age > textClient.world.GetAge() {
				textClient.world.SetAge(age)

				entity, err := textClient.api.GetPlayer(textClient.username)
				if err != nil {
					textClient.Quit(err)
					return
				}

				textClient.player.Entity = entity

				center := textClient.player.Entity.Location.Chunk

				chunks, err := textClient.api.GetChunks(center, constants.ACTIVE_CHUNK_RADIUS)
				if err != nil {
					textClient.Quit(err)
					return
				}

				for i, chunk := range chunks.Chunks {
					point := chunks.Points[i]
					textClient.world.SetChunk(point, chunk)
				}
			}
		}
	}
}

func (textClient *TextClient) displayDiameter() uint {
	return textClient.displayRadius*2 + 1 // 1 is the voxel w/ the player
}

func (textClient *TextClient) Quit(err error) {
	textClient.err = err
	textClient.quitOnce.Do(func() {
		close(textClient.quitq)
	})
}

func (textClient *TextClient) Auth() error {
	textClient.api = client.NewAPI("http://localhost:8080", textClient.username, "murakami")

	if exists, err := textClient.api.UserExists(); err == nil {
		if !exists {
			if err := textClient.api.Signup(); err != nil {
				return err
			}
		}
	} else {
		return err
	}

	if err := textClient.api.Login(); err != nil {
		return err
	}

	entity, err := textClient.api.GetPlayer(textClient.username)
	if err != nil {
		return err
	}

	textClient.player = &types.PlayerDetails{
		Player: &types.Player{
			Name:     textClient.username,
			EntityId: entity.Id,
		},
		Entity: entity,
	}

	return nil
}
