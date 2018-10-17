package main

import (
    "errors"
    "fmt"
    "github.com/felzix/huyilla/types"
    "github.com/loomnetwork/go-loom/plugin"
    contract "github.com/loomnetwork/go-loom/plugin/contractpb"
)


var PLAYERS = []byte("Players")


func (c *Huyilla) GetPlayerList (ctx contract.StaticContext, req *plugin.Request) (*types.PlayerList, error) {
    var players = &types.Players{}
    if err := ctx.Get(PLAYERS, players); err != nil { return nil, err }

    list := types.PlayerList{}
    for name, _ := range players.Players {
        list.Names = append(list.Names, name)
    }
    return &list, nil
}

func (c *Huyilla) GetPlayer (ctx contract.StaticContext, req *types.PlayerName) (*types.Entity, error) {
    var players = &types.Players{}
    if err := ctx.Get(PLAYERS, players); err != nil { return nil, err }
    location := players.Players[req.Name]

    chunk, err := c.getChunk(ctx, location.Chunk)
    if err != nil {
        return nil, err
    }

    var entity *types.Entity
    for i := 0; i < len(chunk.Entities); i++ {
        foundEntity := chunk.Entities[i].Point.X == location.Voxel.X &&
                       chunk.Entities[i].Point.Y == location.Voxel.Y &&
                       chunk.Entities[i].Point.Z == location.Voxel.Z
        if foundEntity {
           entity = chunk.Entities[i].Entity
           break
        }
    }
    if entity == nil {
        return nil, errors.New(fmt.Sprintf(`No such player "%v"`, req.Name))
    }

    return entity, nil
}
