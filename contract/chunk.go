package main

import (
    "fmt"
    "github.com/felzix/huyilla/engine"
    "github.com/felzix/huyilla/types"
    contract "github.com/loomnetwork/go-loom/plugin/contractpb"
)


func (c *Huyilla) GetChunk (ctx contract.StaticContext, req *types.Point) (*types.Chunk, error) {
    return c.getChunk(ctx, req)
}

func (c *Huyilla) GenChunk (ctx contract.Context, req *types.Point) error {
    return c.genChunk(ctx, req)
}

func (c *Huyilla) chunkKey (point *types.Point) []byte {
    return []byte(fmt.Sprintf(`Chunk.%d.%d.%d`, point.X, point.Y, point.Z))
}

func (c *Huyilla) getChunk (ctx contract.StaticContext, point *types.Point) (*types.Chunk, error) {
    key := c.chunkKey(point)
    var chunk types.Chunk
    if err := ctx.Get(key, &chunk); err != nil { return nil, err }
    return &chunk, nil
}

func (c *Huyilla) setChunk (ctx contract.Context, point *types.Point, chunk *types.Chunk) error {
    key := c.chunkKey(point)
    if err := ctx.Set(key, chunk); err != nil { return err }
    return nil
}

func (c *Huyilla) genChunk (ctx contract.Context, point *types.Point) error {
    // iterates over 3D array of voxels.
    // TODO efficient
    chunk := types.Chunk{}
    for y := 0; y < engine.CHUNK_SIZE; y++ {
        for x := 0; x < engine.CHUNK_SIZE; x++ {
            for z := 0; z < engine.CHUNK_SIZE; z++ {
                chunk.Voxels = append(chunk.Voxels, 0x0)  // TODO generate an actual world
            }
        }
    }
    return c.setChunk(ctx, point, &chunk)
}

func (c *Huyilla) addEntityToChunk (ctx contract.Context, entity *types.Entity) error {
    chunk, err := c.getChunk(ctx, entity.Location.Chunk)
    if err != nil { return err }

    chunk.Entities = append(chunk.Entities, entity.Id)

    err = c.setChunk(ctx, entity.Location.Chunk, chunk)

    return nil
}

func (c *Huyilla) removeEntityFromChunk (ctx contract.Context, entity *types.Entity) error {
    chunk, err := c.getChunk(ctx, entity.Location.Chunk)
    if err != nil { return err }

    entities := chunk.Entities
    for i := 0; i < len(entities); i++ {
        id := entities[i]
        if entity.Id == id {
            // idiomatic way of removing a list element in Go
            entities[i] = entities[len(entities) - 1]
            entities = entities[:len(entities) - 1]
            break
        }
    }

    chunk.Entities = append(chunk.Entities, entity.Id)
    return c.setChunk(ctx, entity.Location.Chunk, chunk)
}
