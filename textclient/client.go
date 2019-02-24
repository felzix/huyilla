package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"sync"
	"time"

	react "github.com/felzix/go-curses-react"
	C "github.com/felzix/huyilla/constants"
	"github.com/felzix/huyilla/content"
	"github.com/felzix/huyilla/types"
	"github.com/gdamore/tcell"
)

type Client struct {
	sync.Mutex

	world    *WorldCache
	player   *types.PlayerDetails
	username string

	screen *react.Screen

	quitq    chan struct{}
	quitOnce sync.Once
	err      error

	eventq chan tcell.Event
}

const (
	VIEWMODE_INTRO = iota
	VIEWMODE_GAME
)

func (client *Client) Init() error {
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
	client.screen.Init(root, func(err error) error {
		client.Quit(err)
		return nil
	})

	client.world = &WorldCache{}
	client.world.Init()

	content.PopulateContentNameMaps()

	return nil
}

func (client *Client) Deinit() {
	client.screen.TCellScreen.Fini()
}

func (client *Client) Run() error {
	go client.EventPoller()
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
	client.screen.TCellScreen.PostEvent(iev)

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
	client.screen.HandleKey(e)
	return nil
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

func (client *Client) Quit(err error) {
	client.err = err
	client.quitOnce.Do(func() {
		close(client.quitq)
	})
}

func (client *Client) Auth() error {
	base := "http://localhost:8080"
	res, err := http.Get(base + "/ping")

	fmt.Println(res)
	fmt.Println(err)

	body, err := ioutil.ReadAll(res.Body)

	fmt.Println(body)
	fmt.Println(err)

	// TODO this is blocked by there not yet being a way for client and engine to communicate

	/*
		err := signUp(client.username)
		if err != nil {
			if err.Error() != "rpc error: code = Unknown desc = You are already signed up." {
				return errors.Wrap(err, "Signup error")
			}
		}

		if player, err := logIn(); err == nil {
			client.player = player
		} else if err.Error() == "rpc error: code = Unknown desc = You are already logged in." {
			if addr, err := myAddress(); err != nil {
				return errors.Wrap(err, "MyAddress error")
			} else if player, err := getPlayer(addr); err != nil {
				return errors.Wrap(err, "GetPlayer error")
			} else {
				client.player = player
			}
		} else {
			return errors.Wrap(err, "Login error")
		}
	*/
	return nil
}

func MakeApp() *react.ReactElement {
	root := &react.ReactElement{
		State: react.State{
			"mode": VIEWMODE_INTRO,
		},
		DrawFn: func(r *react.ReactElement, maxWidth, maxHeight int) (*react.DrawResult, error) {
			client := r.Props["client"].(*Client)
			mode := r.State["mode"].(int)

			var element *react.ReactElement
			var props react.Properties
			switch mode {
			case VIEWMODE_INTRO:
				element = Intro()
				props = react.Properties{
					"client": client,
					"nextMode": func() {
						r.State["mode"] = VIEWMODE_GAME
					},
				}
			case VIEWMODE_GAME:
				element = GameBoard()
				props = react.Properties{
					"client": client,
				}
			}

			result := react.DrawResult{
				Elements: []react.Child{
					*react.NewChild(element, string(mode), maxWidth, maxHeight, props),
				}}
			return &result, nil
		},
	}

	return root
}

func Intro() *react.ReactElement {
	return &react.ReactElement{
		Type: "Intro",
		DrawFn: func(r *react.ReactElement, maxWidth, maxHeight int) (*react.DrawResult, error) {
			client := r.Props["client"].(*Client)
			nextMode := r.Props["nextMode"].(func())

			child := react.NewChild(react.HorizontalLayout(), "", maxWidth, maxHeight, react.Properties{
				"children": []*react.Child{
					react.ManagedChild(react.Label(), "hello", react.Properties{
						"label": "Hello!",
					}),
					react.ManagedChild(react.Label(), "blank", react.Properties{
						"label": "",
					}),
					react.ManagedChild(react.TextEntry(), "", react.Properties{
						"label": "Enter username",
						"whenFinished": func(username string) error {
							client.username = username
							nextMode()
							return nil
						},
						// TODO when TextEntry can do validation, reject empty or taken username
					}),
				},
			})
			result := react.DrawResult{
				Elements: []react.Child{*child},
			}
			return &result, nil
		},
	}
}

// TODO this should be polling the engine, once that's possible
func GameBoard() *react.ReactElement {
	return &react.ReactElement{
		Type: "GameBoard",
		DrawFn: func(r *react.ReactElement, maxWidth, maxHeight int) (*react.DrawResult, error) {
			client := r.Props["client"].(*Client)

			// there is no client.player
			child := react.NewChild(react.HorizontalLayout(), "", maxWidth, maxHeight, react.Properties{
				"children": []*react.Child{
					react.ManagedChild(react.Label(), "debug-bar", react.Properties{
						"label": fmt.Sprintf("%d", client.world.age),
					}),
					react.ManagedChild(react.Label(), "blank", react.Properties{
						"label": "",
					}),
					react.ManagedChild(Tiles(), "", react.Properties{
						"client":   client,
						"absPoint": client.player.Entity.Location,
					}),
				},
			})
			result := react.DrawResult{
				Elements: []react.Child{*child},
			}
			return &result, nil
		},
	}
}

func Tiles() *react.ReactElement {
	return &react.ReactElement{
		Type: "Tiles",
		DrawFn: func(r *react.ReactElement, maxWidth, maxHeight int) (*react.DrawResult, error) {
			client := r.Props["client"].(*Client)
			absPoint := r.Props["absPoint"].(*types.AbsolutePoint)

			chunk := client.world.chunks[*absPoint.Chunk]
			zLevel := int(absPoint.Voxel.Z)

			width := C.CHUNK_SIZE
			if width > maxWidth {
				width = maxWidth
			}
			height := C.CHUNK_SIZE
			if height > maxHeight {
				height = maxHeight
			}

			result := react.DrawResult{
				Region: react.NewRegion(0, 0, maxWidth, maxHeight),
			}

			for y := 0; y < height; y++ {
				for x := 0; x < width; x++ {
					if chunk == nil {
						result.Region.Cells[x][y] = react.Cell{
							R:     ' ',
							Style: tcell.StyleDefault.Background(tcell.ColorDarkGray),
						}
					} else {
						index := (x * C.CHUNK_SIZE * C.CHUNK_SIZE) + (y * C.CHUNK_SIZE) + zLevel
						ch := voxelToRune(chunk.Voxels[index])
						result.Region.Cells[x][y] = react.Cell{
							R:     ch,
							Style: tcell.StyleDefault,
						}
					}
				}
			}
			return &result, nil
		},
	}
}

func voxelToRune(voxel uint64) rune {
	voxelType := voxel & 0xFFFF

	switch voxelType {
	case content.VOXEL["air"]:
		return ' '
	case content.VOXEL["barren_earth"]:
		return '.'
	case content.VOXEL["barren_grass"]:
		return ','
	case content.VOXEL["water"]:
		return '~'
	default:
		return rune(0)
	}
}
