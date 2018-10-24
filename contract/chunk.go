package main

import (
    "fmt"
    C "github.com/felzix/huyilla/constants"
    "github.com/felzix/huyilla/types"
    contract "github.com/loomnetwork/go-loom/plugin/contractpb"
    "github.com/mitchellh/hashstructure"
    "math/rand"
)


func (c *Huyilla) GetChunk (ctx contract.StaticContext, req *types.Point) (*types.Chunk, error) {
    return c.getChunk(ctx, req)
}

func (c *Huyilla) GenChunk (ctx contract.Context, p *types.Point) error {
    return c.genChunk(ctx, C.SEED, p)
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

func (c *Huyilla) getChunkGuaranteed (ctx contract.Context, point *types.Point) (*types.Chunk, error) {
    key := c.chunkKey(point)
    var chunk types.Chunk

    if !ctx.Has(key) {
        if err := c.genChunk(ctx, C.SEED, point); err != nil {
            return nil, err
        }
    }

    if err := ctx.Get(key, &chunk); err != nil {
        return nil, err
    } else {
        return &chunk, nil
    }
}

func (c *Huyilla) setChunk (ctx contract.Context, point *types.Point, chunk *types.Chunk) error {
    key := c.chunkKey(point)
    if err := ctx.Set(key, chunk); err != nil { return err }
    return nil
}

func (c *Huyilla) addEntityToChunk (ctx contract.Context, entity *types.Entity) error {
    chunk, err := c.getChunkGuaranteed(ctx, entity.Location.Chunk)
    if err != nil { return err }

    chunk.Entities = append(chunk.Entities, entity.Id)

    err = c.setChunk(ctx, entity.Location.Chunk, chunk)

    return nil
}

func (c *Huyilla) removeEntityFromChunk (ctx contract.Context, entity *types.Entity) error {
    chunk, err := c.getChunk(ctx, entity.Location.Chunk)

    if err != nil {
        if err.Error() == "not found" {
            return nil  // chunk doesn't exist anyway so it need not be changed
        } else {
            return err // something else went wrong
        }
    }

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


func (c *Huyilla) genChunk (ctx contract.Context, worldSeed uint64, p *types.Point) error {
    chunkSeed, _ := hashstructure.Hash(p, nil)
    seed := int64(worldSeed * chunkSeed)

    chunk := types.Chunk{Voxels: make([]uint64, C.CHUNK_LENGTH)}
    var x, y, z int64
    for x = 0; x < C.CHUNK_SIZE; x++ {
        for y = 0; y < C.CHUNK_SIZE; y++ {
            for z = 0; z < C.CHUNK_SIZE; z++ {
                rand.Seed(seed)  // so voxels can use randomness
                index := (x * C.CHUNK_SIZE * C.CHUNK_SIZE) + (y * C.CHUNK_SIZE) + z
                location := &types.AbsolutePoint{
                    Chunk: p,
                    Voxel: &types.Point{x, y, z},
                }
                chunk.Voxels[index] = genVoxel(location)
            }
        }
    }
    return c.setChunk(ctx, p, &chunk)
}

func genVoxel (p *types.AbsolutePoint) uint64 {
    v := VOXEL

    if p.Chunk.Z < 0 {
        return v["dirt"]
    }

    if p.Chunk.Z > 0 {
        return v["air"]
    }

    center := randomPoint()
    center.Z = 0

    d := distance(p.Voxel, center)
    if p.Voxel.Z == center.Z && d <= float64(3) {
        return v["water"]
    }
    if p.Voxel.Z == 0 {
        return v["barren_earth"]
    }
    return v["air"]
}