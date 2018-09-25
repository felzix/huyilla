package engine

import (
    "fmt"
    "testing"
)

func TestChunks(t *testing.T) {
    world := NewWorld(getContent(t), NewPondWorldGenerator(12, 2))
    world.GenerateChunk(Point{0, 0, 0})
    chunk := world.GetChunk(Point{0,0,0})

    // helpful if there's a failure
    for y := 0; y < CHUNK_SIZE; y++ {
        for x := 0; x < CHUNK_SIZE; x++ {
            voxel := chunk.Get(Point{x, y, 0})
            fmt.Print(voxel.Type)
        }
        fmt.Println()
    }

    voxel := chunk.Get(Point{9, 1, 0})
    if vt := voxel.Type; vt != world.Content.V["water"] {
        t.Errorf(`Generator should have made water("%d") but instead made "%d"`,
            world.Content.V["water"], vt)
    }

}
