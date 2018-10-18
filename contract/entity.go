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

func (c *Huyilla) makeEntityKey (id int64) []byte {
    return []byte(fmt.Sprintf(`Entity.%d`, id))
}

func (c *Huyilla) newEntity (typeInt uint32, playerName string) (*types.Entity, error) {
    entity := types.Entity{
        Id: c.genRandomId(),
        Type: typeInt,
    }
    if playerName == "" {
        entity.Control = types.Entity_NPC
    } else {
        entity.Control = types.Entity_PLAYER
        entity.PlayerName = playerName
    }

    return &entity, nil
}


func (c *Huyilla) setEntity (ctx contract.Context, entity *types.Entity) error {
    return ctx.Set(c.makeEntityKey(entity.Id), entity)
}

func (c *Huyilla) getEntity (ctx contract.StaticContext, id int64) (*types.Entity, error) {
    var entity types.Entity
    err := ctx.Get(c.makeEntityKey(id), &entity)
    if err != nil {return nil, err}
    return &entity, nil
}

func (c *Huyilla) genRandomId () int64 {
    return rand.Int63()
}