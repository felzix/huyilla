package main

import (
    "github.com/felzix/huyilla/types"
    "github.com/loomnetwork/go-loom/plugin"
    contract "github.com/loomnetwork/go-loom/plugin/contractpb"
)



func (c *Huyilla) Tick (ctx contract.Context, req *plugin.Request) error {
    actions, err := c.getActions(ctx)
    if err != nil { return err }

    players, err := c.getActivePlayers(ctx)
    if err != nil { return err }

    activeChunks := make(map[types.Point]bool, len(players) * ACTIVE_CHUNK_CUBE)
    for i := 0; i < len(players); i++ {
        player := players[i]
        loc := player.Entity.Location.Chunk
        for x := loc.X - ACTIVE_CHUNK_RADIUS; x < loc.X + ACTIVE_CHUNK_RADIUS; x++ {
            for y := loc.Y - ACTIVE_CHUNK_RADIUS; y < loc.Y + ACTIVE_CHUNK_RADIUS; y++ {
                for z := loc.Z - ACTIVE_CHUNK_RADIUS; z < loc.Z + ACTIVE_CHUNK_RADIUS; z++ {
                    activeChunks[*newPoint(x, y, z)] = true
                }
            }
        }
    }

    vitalizedVoxels := make([]types.Point, PASSIVE_VITALITY)
    for i := 0; i < PASSIVE_VITALITY; i++ {
        vitalizedVoxels[i] = *randomPoint()
    }

    for p, _ := range activeChunks {
        chunk, err := c.getChunk(ctx, &p)
        if err != nil { return err }

        for i := 0; i < len(chunk.ActiveVoxels); i++ {
            err := c.voxelPhysics(ctx, chunk, &types.AbsolutePoint{Chunk: &p, Voxel: chunk.ActiveVoxels[i]} )
            if err != nil { return err }
        }

        for i := 0; i < PASSIVE_VITALITY; i++ {
            err := c.voxelPhysics(ctx, chunk, &types.AbsolutePoint{Chunk: &p, Voxel: &vitalizedVoxels[i]})
            if err != nil { return err }
        }

        for i := 0; i < len(chunk.Entities); i++ {
            var entity types.Entity
            err := ctx.Get(c.entityKey(chunk.Entities[i]), &entity)
            if err != nil { return err }

            err = c.entityPhysics(ctx, chunk, &entity)
            if err != nil { return err }
        }
    }

    for i := 0; i < len(actions.Actions); i++ {
        // TODO
        // if valid then apply then emit success event
        // else, ignore then emit failure event
    }

    // reset actions queue
    return ctx.Set(ACTIONS, &types.Actions{})
}

func (c *Huyilla) voxelPhysics (ctx contract.Context, chunk *types.Chunk, location *types.AbsolutePoint) error {
    return nil
}

func (c *Huyilla) entityPhysics (ctx contract.Context, chunk *types.Chunk, entity *types.Entity) error {
    return nil
}
