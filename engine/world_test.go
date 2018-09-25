package engine

import (
    "fmt"
    "testing"
)

func TestChunkGeneration(t *testing.T) {
    world := NewWorld(getContent(t), NewPondWorldGenerator(12, 2))
    chunk := world.GetChunk(Point{0,0,0})

    c2 := world.GetChunk(Point{0, 0, 0})
    if chunk != c2 {
        t.Error("Getting a chunk twice should yield exactly the same chunk but didn't")
    }

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
