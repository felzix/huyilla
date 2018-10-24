package main

import (
    "encoding/json"
    "fmt"
    C "github.com/felzix/huyilla/constants"
    "github.com/felzix/huyilla/types"
    "github.com/loomnetwork/go-loom/plugin"
    contract "github.com/loomnetwork/go-loom/plugin/contractpb"
)



func (c *Huyilla) Tick (ctx contract.Context, req *plugin.Request) error {
    players, err := c.getActivePlayers(ctx)
    if err != nil { return err }

    activeChunks := make(map[types.Point]*types.Chunk, len(players) * C.ACTIVE_CHUNK_CUBE)
    for i := 0; i < len(players); i++ {
        player := players[i]
        loc := player.Entity.Location.Chunk
        for x := loc.X - C.ACTIVE_CHUNK_RADIUS; x < 1 + loc.X + C.ACTIVE_CHUNK_RADIUS; x++ {
            for y := loc.Y - C.ACTIVE_CHUNK_RADIUS; y < 1 + loc.Y + C.ACTIVE_CHUNK_RADIUS; y++ {
                for z := loc.Z - C.ACTIVE_CHUNK_RADIUS; z < 1 + loc.Z + C.ACTIVE_CHUNK_RADIUS; z++ {
                    point := newPoint(x, y, z)
                    chunk, err := c.getChunkGuaranteed(ctx, point)
                    if err != nil { return err }
                    activeChunks[*point] = chunk
                }
            }
        }
    }

    vitalizedVoxels := make([]types.Point, C.PASSIVE_VITALITY)
    for i := 0; i < C.PASSIVE_VITALITY; i++ {
        vitalizedVoxels[i] = *randomPoint()
    }

    for p, chunk := range activeChunks {
        if err != nil { return err }

        for i := 0; i < len(chunk.ActiveVoxels); i++ {
            err := c.voxelPhysics(ctx, chunk, &types.AbsolutePoint{Chunk: &p, Voxel: chunk.ActiveVoxels[i]} )
            if err != nil { return err }
        }

        for i := 0; i < C.PASSIVE_VITALITY; i++ {
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

    actions, err := c.getActions(ctx)
    if err != nil { return err }

    for i := 0; i < len(actions.Actions); i++ {
        action := actions.Actions[i]

        var fn func(contract.Context, *types.Action) (bool, error)
        var actionName string
        switch a := action.Action.(type) {
        case *types.Action_Move:
            fn = c.move
            actionName = "move"
        default:
            // only log error - if the action is broken then don't block the engine
            ctx.Logger().Error(fmt.Sprintf("Invalid action %v", a))
            continue
        }

        success, err := fn(ctx, action)
        if err != nil {
            // only log error - if the action is broken then don't block the engine
            ctx.Logger().Error(err.Error())
            continue
        }

        emitMsg := struct {
            Action  string
            Addr    []byte
            Success bool
        }{actionName, c.thisUser(ctx), success}
        emitMsgJSON, err := json.Marshal(emitMsg)
        if err != nil {
            // only log error - if the action is broken then don't block the engine
            ctx.Logger().Error(err.Error())
            continue
        }
        ctx.EmitTopics(emitMsgJSON, "huyilla:action:" + string(emitMsg.Addr))
    }

    // clear actions queue
    ctx.Delete(ACTIONS)

    // save chunks
    for p, chunk := range activeChunks {
        c.setChunk(ctx, &p, chunk)
    }

    // advance age by one tick
    if _, err := c.incrementAge(ctx); err != nil {
        return err
    }

    return nil
}

func (c *Huyilla) voxelPhysics (ctx contract.Context, chunk *types.Chunk, location *types.AbsolutePoint) error {
    return nil
}

func (c *Huyilla) entityPhysics (ctx contract.Context, chunk *types.Chunk, entity *types.Entity) error {
    return nil
}
