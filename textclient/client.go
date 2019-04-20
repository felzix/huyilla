package main

import (
	"sync"
	"time"

	react "github.com/felzix/go-curses-react"
	"github.com/felzix/huyilla/content"
	"github.com/felzix/huyilla/types"
	"github.com/gdamore/tcell"
)

type Client struct {
	sync.Mutex

	world    *WorldCache
	player   *types.PlayerDetails
	username string
	api      *API

	screen *react.Screen
	viewDepthDelta int
	displayRadius uint

	quitq    chan struct{}
	quitOnce sync.Once
	err      error

	eventq chan tcell.Event
}

func (client *Client) Init() error {
	client.displayRadius = 15

	if screen, err := react.NewScreen(); err == nil {
		client.screen = screen
	} else {
		return err
	}

	client.quitq = make(chan struct{})
	client.eventq = make(chan tcell.Event)

	root := MakeApp()
	root.Props = react.Properties{
		"client": client,
	}
	if err := client.screen.Init(root, func(err error) error {
		client.Quit(err)
		return nil
	}); err != nil {
		return err
	}

	client.world = &WorldCache{}
	client.world.Init()

	content.PopulateContentNameMaps()

	return nil
}

func (client *Client) Deinit() {
	client.screen.TCellScreen.Fini()
}

func (client *Client) Run() error {
	go client.EventPoller() // gets user input
	go client.EnginePoller() // gets world state from the engine

loop:
	for {
		if err := client.Draw(); err != nil {
			client.Quit(err)
			break loop
		}
		select {
		case <-client.quitq:
			break loop
		// draw-input loop runs no faster than once every 10ms
		case <-time.After(time.Millisecond * 10):
		case ev := <-client.eventq:
			if err := client.HandleEvent(ev); err != nil {
				client.Quit(err)
				break loop
			}
		}
	}

	// Inject a wakeup interrupt
	iev := tcell.NewEventInterrupt(nil)
	if err := client.screen.TCellScreen.PostEvent(iev); err != nil {
		return err
	}

	return client.err
}

func (client *Client) HandleEvent(e tcell.Event) error {
	switch e := e.(type) {
	case *tcell.EventResize:
		client.handleResize(e)
	case *tcell.EventKey:
		return client.handleKey(e)
	case *tcell.EventMouse:
		return client.handleMouse(e)
	}

	return nil
}

func (client *Client) handleResize(e *tcell.EventResize) {
	client.screen.Resize()
	client.screen.TCellScreen.Sync() // visually jarring but needed after a resize
}

func (client *Client) handleKey(e *tcell.EventKey) error {
	client.Lock()
	defer client.Unlock()
	return client.screen.HandleKey(e)
}

func (client *Client) handleMouse(e *tcell.EventMouse) error {
	// x, y := e.Position()
	// fmt.Printf("(%d,%d)", x, y)
	return nil
}

func (client *Client) Draw() error {
	client.Lock()
	defer client.Unlock()

	err := client.screen.Draw()
	if err == nil {
		client.screen.TCellScreen.Show()
	}
	return err
}

func (client *Client) EventPoller() {
	for {
		select {
		case <-client.quitq:
			return
		default:
		}

		e := client.screen.TCellScreen.PollEvent()
		if e == nil {
			return
		}

		select {
		case <-client.quitq:
			return
		case client.eventq <- e:
		}
	}
}

func (client *Client) EnginePoller() {
	for {
		select {
		case <- client.quitq:
			return
		case <- time.After(time.Millisecond * 1000): // poll engine only so often
		default:
		}

		if client.api == nil || client.player == nil {
			continue // user is still entering in their information
		}

		// ignores error because getting world age is eqiuvalent to querying the readiness of the server
		age, _ := client.api.GetWorldAge()

		client.world.age = age

		if client.world.age == 0 {
			continue // world is not loaded
		}

		if client.player.Entity == nil {
			entity, err := client.api.GetPlayer(client.username)
			if err != nil {
				client.Quit(err)
				return
			}

			client.player.Entity = entity
		}

		centerChunk := client.player.Entity.Location.Chunk
		chunk, err := client.api.GetChunk(centerChunk)
		if err != nil {
			client.Quit(err)
			return
		}

		client.SetChunk(centerChunk, chunk)
	}
}

func (client *Client) displayDiameter() uint {
	return client.displayRadius * 2 + 1 // 1 is the voxel w/ the player
}

func (client *Client) Quit(err error) {
	client.err = err
	client.quitOnce.Do(func() {
		close(client.quitq)
	})
}

func (client *Client) Auth() error {
	client.api = NewAPI("http://localhost:8080", client.username, "murakami")

	if exists, err := client.api.UserExists(); err == nil {
		if !exists {
			if err := client.api.Signup(); err != nil {
				return err
			}
		}
	} else {
		return err
	}

	if err := client.api.Login(); err != nil {
		return err
	}

	entity, err := client.api.GetPlayer(client.username)
	if err != nil {
		return err
	}

	client.player = &types.PlayerDetails{
		Player: &types.Player{
			Name: client.username,
			EntityId: entity.Id,
		},
		Entity: entity,
	}

	return nil
}
