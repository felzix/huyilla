package main

import (
    "encoding/base64"
    "fmt"
    C "github.com/felzix/huyilla/constants"
    "github.com/felzix/huyilla/types"
    "github.com/gdamore/tcell"
    "github.com/gdamore/tcell/views"
    "github.com/pkg/errors"
    "golang.org/x/crypto/ed25519"
    "io/ioutil"
    "log"
    "sync"
    "time"
)


type Client struct {
    sync.Mutex

    world *WorldCache

    screen tcell.Screen
    view *views.ViewPort
    debugView *views.ViewPort

    quitq chan struct{}
    quitOnce sync.Once
    err error

    eventq chan tcell.Event
}


func (client *Client) Init () error {
    client.world = &WorldCache{}
    client.world.Init()

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

    client.view = views.NewViewPort(client.screen, 0, 2, C.CHUNK_SIZE, C.CHUNK_SIZE)
    client.debugView = views.NewViewPort(client.screen, 0, 0, 80, 1)

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
        width, height := e.Size()

        client.view.Resize(0, 2, width, height)
        client.debugView.Resize(0, 0, width, 1)
    case *tcell.EventKey:
        switch e.Key(){
        case tcell.KeyEsc:
            client.Quit(nil)
        }
    case *tcell.EventMouse:
        // x, y := e.Position()
        // fmt.Printf("(%d,%d)", x, y)
    }

    return nil
}

func (client *Client) Draw() {
    client.Lock()

    client.view.Clear()
    for y := 0; y < C.CHUNK_SIZE; y++ {
        for x := 0; x < C.CHUNK_SIZE; x++ {
            client.view.SetContent(x, y, '.', nil, tcell.StyleDefault)
        }
    }

    client.debugView.Clear()
    age := fmt.Sprintf("%d", client.world.age)
    for i := 0; i < len(age); i++ {
        client.debugView.SetContent(i, 0, rune(age[i]), nil, tcell.StyleDefault)
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



func generateKey(privFile string) error {
    _, priv, err := ed25519.GenerateKey(nil)
    if err != nil {
        return errors.Wrapf(err, "Error generating key pair")
    }
    data := base64.StdEncoding.EncodeToString(priv)
    if err := ioutil.WriteFile(privFile, []byte(data), 0664); err != nil {
        return errors.Wrapf(err, "Unable to write private key")
    }
    return nil
}

func getAge () (uint64, error) {
    var age types.Age
    if err := StaticCallContract("GetAge", &types.Nothing{}, &age); err != nil {
        return 0, err
    }

    return age.Ticks, nil
}

func getConfig () (map[string]interface{}, error) {
    var config types.Config

    if err := StaticCallContract("GetConfig", &types.Nothing{}, &config); err != nil {
        return nil, err
    }

    native := make(map[string]interface{})
    for k, v := range config.Options.Map {
        switch value := v.Value.(type) {
        case *types.Primitive_Int: native[k] = value.Int
        case *types.Primitive_Bool: native[k] = value.Bool
        case *types.Primitive_String_: native[k] = value.String_
        case *types.Primitive_Float: native[k] = value.Float
        default: native[k] = nil
        }
    }

    return native, nil
}

func getChunk (point *types.Point) (*types.Chunk, error) {
    var chunk types.Chunk

    if err := StaticCallContract("GetChunk", point, &chunk); err != nil {
        return nil, err
    }

    log.Print(chunk.Voxels)

    return &chunk, nil
}
