package main

import (
    "fmt"
    "github.com/felzix/huyilla/types"
    contract "github.com/loomnetwork/go-loom/plugin/contractpb"
    "math/rand"
)


func (c *Huyilla) GetEntity (ctx contract.StaticContext, id *types.EntityId) (*types.Entity, error) {
    return c.getEntity(ctx, id.Id)
}

func (c *Huyilla) entityKey (id int64) []byte {
    return []byte(fmt.Sprintf(`Entity.%d`, id))
}

func (c *Huyilla) newEntity (ctx contract.StaticContext, typeInt uint32, playerName string, location *types.AbsolutePoint) *types.Entity {
    entity := types.Entity{
        Id: c.genUniqueEntityId(ctx),
        Type: typeInt,
        Location: location,
    }
    if playerName == "" {
        entity.Control = types.Entity_NPC
    } else {
        entity.Control = types.Entity_PLAYER
        entity.PlayerName = playerName
    }

    return &entity
}


func (c *Huyilla) setEntity (ctx contract.Context, entity *types.Entity) error {
    oldEntity, _ := c.getEntity(ctx, entity.Id)

    // sync chunks' entity lists
    if oldEntity == nil {
        err := c.addEntityToChunk(ctx, entity)
        if err != nil { return err }
    } else if pointEquals(oldEntity.Location.Chunk, entity.Location.Chunk) {
        if err := c.removeEntityFromChunk(ctx, oldEntity); err != nil {
            return err
        }
        if err := c.addEntityToChunk(ctx, entity); err != nil {
            return err
        }
    }

    return ctx.Set(c.entityKey(entity.Id), entity)
}

func (c *Huyilla) getEntity (ctx contract.StaticContext, id int64) (*types.Entity, error) {
    var entity types.Entity
    err := ctx.Get(c.entityKey(id), &entity)
    if err != nil {return nil, err}
    return &entity, nil
}

func (c *Huyilla) genUniqueEntityId (ctx contract.StaticContext) int64 {
    var id int64
    for true {
        id = rand.Int63()
        if !ctx.Has(c.entityKey(id)) {
            break
        }
    }
    return id
}