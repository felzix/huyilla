package main

import (
    "fmt"
    C "github.com/felzix/huyilla/constants"
    "github.com/gdamore/tcell"
    "github.com/gdamore/tcell/views"
    "sync"
    "time"
)


type Client struct {
    sync.Mutex

    world *WorldCache
    username string
    viewMode int

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
    client.screen.EnableMouse()  // TODO do I want this?

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
loop:
    for {
        client.Draw()
        select {
        case <-client.quitq:
            break loop
        // draw-input loop runs no faster than once every 10ms
        case <-time.After(time.Millisecond * 10):
        case ev := <-client.eventq:
            client.HandleEvent(ev)
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
                client.viewMode = VIEWMODE_GAME
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

func (client *Client) Draw() {
    client.Lock()

    switch client.viewMode {
    case VIEWMODE_INTRO:
        client.introView.Clear()
        drawString(client.introView, 0, 0, "Hello!")
        drawString(client.introView, 0, 2, "Enter username: " + client.username)
    case VIEWMODE_GAME:
        client.chunkView.Clear()
        for y := 0; y < C.CHUNK_SIZE; y++ {
            for x := 0; x < C.CHUNK_SIZE; x++ {
                client.chunkView.SetContent(x, y, '.', nil, tcell.StyleDefault)
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
}

func (client *Client) Updater() {
    for {
        select {
        case <-client.quitq:
            return
        // tick-query loop runs no faster than once every 10ms
        case <-time.After(time.Millisecond * 50):
            client.Lock()

            // TODO tick engine then query its state
            if age, err := getAge(); err == nil {
                client.world.age = age
            } else {
                client.Quit(err)
            }

            client.Unlock()
        }
    }
}

func (client *Client) EventPoller() {
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



func drawString(view *views.ViewPort, x, y int, s string) {
    for i := 0; i < len(s); i++ {
        view.SetContent(x + i, y, rune(s[i]), nil, tcell.StyleDefault)
    }
}
