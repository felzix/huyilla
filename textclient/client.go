package main

import (
    "fmt"
    C "github.com/felzix/huyilla/constants"
    "github.com/felzix/huyilla/content"
    "github.com/felzix/huyilla/types"
    "github.com/gdamore/tcell"
    "github.com/gdamore/tcell/views"
    "github.com/loomnetwork/go-loom"
    "github.com/loomnetwork/go-loom/auth"
    "github.com/pkg/errors"
    "sync"
    "time"
)


type Client struct {
    sync.Mutex

    world *WorldCache
    player *types.PlayerDetails
    username string
    viewMode int
    signer *auth.Ed25519Signer

    screen tcell.Screen
    introView *views.ViewPort
    chunkView *views.ViewPort
    debugView *views.ViewPort

    quitq chan struct{}
    quitOnce sync.Once
    err error

    eventq chan tcell.Event
}

const (
    VIEWMODE_INTRO = 0
    VIEWMODE_GAME = 1
)


func (client *Client) Init () error {
    client.world = &WorldCache{}
    client.world.Init()

    client.viewMode = VIEWMODE_INTRO

    if signer, err := MakeSigner(); err != nil {
        return err
    } else {
        client.signer = signer
    }

    if screen, err := tcell.NewScreen(); err != nil {
        return err
    } else if err = screen.Init(); err != nil {
        return err
    } else {
        client.screen = screen
    }
    client.screen.SetStyle(tcell.StyleDefault.
        Background(tcell.ColorBlack).
        Foreground(tcell.ColorWhite))
    // client.screen.EnableMouse()  // TODO do I want this?

    width, height := client.screen.Size()
    client.introView = views.NewViewPort(client.screen, 0, 0, width, height)
    client.chunkView = views.NewViewPort(client.screen, 0, 2, width, height - 2)
    client.debugView = views.NewViewPort(client.screen, 0, 0, width, 1)

    client.quitq = make(chan struct{})
    client.eventq = make(chan tcell.Event)

    return nil
}

func (client *Client) Deinit () {
    client.screen.Fini()
}

func (client *Client) Run () error {
    go client.EventPoller()
    go client.Updater()
    go client.Ticker()
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
    client.screen.PostEvent(iev)

    return client.err
}

func (client *Client) HandleEvent (e tcell.Event) error {
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

func (client *Client) handleResize (e *tcell.EventResize) {
    width, height := e.Size()
    client.introView.Resize(0, 0, width, height)
    client.chunkView.Resize(0, 2, width, height - 2)
    client.debugView.Resize(0, 0, width, 1)
}

func (client *Client) handleKey (e *tcell.EventKey) error {
    if e.Key() == tcell.KeyEsc {
        client.Quit(nil)
        return nil
    }

    switch client.viewMode {
    case VIEWMODE_INTRO:
        if e.Key() == tcell.KeyEnter {
            if len(client.username) > 0 {
                client.Lock()
                client.viewMode = VIEWMODE_GAME
                err := client.Auth()
                client.Unlock()
                return err
            }
        } else if e.Rune() != 0 {
            client.username += string(e.Rune())
        }
    case VIEWMODE_GAME:
        switch e.Key() {
        case tcell.KeyUp:
            // TODO issue move command to server for player entity to move up/north, depending on terrain
        }

        switch e.Rune() {
        case 'q':
            client.Quit(nil)
        }
    }
    return nil
}

func (client *Client) handleMouse (e *tcell.EventMouse) error {
    // x, y := e.Position()
    // fmt.Printf("(%d,%d)", x, y)
    return nil
}

func (client *Client) Draw() error {
    client.Lock()

    switch client.viewMode {
    case VIEWMODE_INTRO:
        client.introView.Clear()
        drawString(client.introView, 0, 0, "Hello!")
        drawString(client.introView, 0, 2, "Enter username: " + client.username)
    case VIEWMODE_GAME:
        point := client.player.Entity.Location.Chunk
        zLevel := int(client.player.Entity.Location.Voxel.Z)
        chunk := client.world.chunks[*point]
        client.chunkView.Clear()
        for y := 0; y < C.CHUNK_SIZE; y++ {
            for x := 0; x < C.CHUNK_SIZE; x++ {
                if chunk == nil {
                    style := tcell.StyleDefault.Background(tcell.ColorDarkGray)
                    client.chunkView.SetContent(int(x), int(y), ' ', nil, style)
                } else {
                    index := (x * C.CHUNK_SIZE * C.CHUNK_SIZE) + (y * C.CHUNK_SIZE) + zLevel
                    ch := voxelToRune(chunk.Voxels[index])
                    client.chunkView.SetContent(x, y, ch, nil, tcell.StyleDefault)
                }
            }
        }

        client.debugView.Clear()
        age := fmt.Sprintf("%d", client.world.age)
        for i := 0; i < len(age); i++ {
            client.debugView.SetContent(i, 0, rune(age[i]), nil, tcell.StyleDefault)
        }
    }

    client.screen.Show()
    client.Unlock()

    return nil
}

func (client *Client) Updater () {
    for {
        select {
        case <-client.quitq:
            return
        // query loop runs no faster than once every 500ms
        case <-time.After(time.Millisecond * 500):
            if client.viewMode == VIEWMODE_GAME {
                client.Lock()

                if client.player != nil {
                    point := client.player.Entity.Location.Chunk
                    if chunk, err := getChunk(point); err == nil {
                        client.world.chunks[*point] = chunk
                    } else {
                        client.Quit(errors.Wrap(err, "GetChunk error"))
                    }
                }

                client.Unlock()
            }
        }
    }
}

func (client *Client) Ticker () {
    for {
        select {
        case <-client.quitq:
            return
        case <-time.After(time.Millisecond * 50):
            if client.viewMode == VIEWMODE_GAME {
                if age, err := getAge(); err == nil {
                    if age > client.world.age {
                        client.Lock()
                        client.world.age = age
                        client.Unlock()

                        if err := tick(); err != nil {
                            client.Quit(errors.Wrap(err, "Tick error"))
                        }
                    }
                } else {
                    client.Quit(errors.Wrap(err, "GetAge error"))
                }
            }
        }
    }
}

func (client *Client) EventPoller () {
    for {
        select {
        case <-client.quitq:
            return
        default:
        }

        e := client.screen.PollEvent()
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

func (client *Client) Quit (err error) {
    client.err = err
    client.quitOnce.Do(func () {
        close(client.quitq)
    })
}

func (client *Client) Auth () error {
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

    return nil
}

func (client *Client) playerAddr () string {
    return loom.Address{"", client.signer.PublicKey()}.Local.String()
}

func drawString(view *views.ViewPort, x, y int, s string) {
    for i := 0; i < len(s); i++ {
        view.SetContent(x + i, y, rune(s[i]), nil, tcell.StyleDefault)
    }
}

func voxelToRune (voxel uint64) rune {
    voxelType := voxel & 0xFFFF

    charMap := map[string]rune {
        "air": ' ',
        "barren_earth": '.',
        "barren_grass": ',',
        "water": '~',
    }

    typeToRune := make(map[uint64]rune, len(charMap))
    for name, rune := range charMap {
        typeToRune[content.VOXEL[name]] = rune
    }

    return typeToRune[voxelType]
}
