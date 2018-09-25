package engine

import (
    "fmt"
    "testing"
)

func TestChunks(t *testing.T) {
    world := MakeWorld(getContent(t))

    fmt.Println(world.Chunks)
    world.GenerateChunk(Point{0, 0, 0})

    // t.Error()

    // c := world.GetChunk(Point{0, 0, 0})
    // fmt.Println(c == nil)
    // t.Errorf("%v", c)
}
