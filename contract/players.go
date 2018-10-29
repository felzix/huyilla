package main

import (
    "github.com/felzix/huyilla/types"
    "github.com/loomnetwork/go-loom/plugin"
    contract "github.com/loomnetwork/go-loom/plugin/contractpb"
    "github.com/pkg/errors"
)


var PLAYERS = []byte("Players")


func (c *Huyilla) GetPlayerList (ctx contract.StaticContext, req *plugin.Request) (*types.PlayerList, error) {
    players, err := c.getPlayers(ctx)
    if err != nil { return nil, err }

    list := types.PlayerList{}
    for _, player := range players.Players {
        list.Names = append(list.Names, player.Name)
    }
    return &list, nil
}

func (c *Huyilla) GetPlayer (ctx contract.StaticContext, req *types.Address) (*types.PlayerDetails, error) {
    return c.getPlayer(ctx, req.Addr)
}

func (c *Huyilla) getPlayer (ctx contract.StaticContext, addr string) (*types.PlayerDetails, error) {
    players, err := c.getPlayers(ctx)
    if err != nil { return nil, err }

    player := players.Players[addr]

    if player == nil {
        return nil, errors.New("No such player " + addr)
    }

    entity, _ := c.getEntity(ctx, player.EntityId)
    // NOTE: entity can be nil

    return &types.PlayerDetails{Player: player, Entity: entity}, nil
}

func (c *Huyilla) getPlayers (ctx contract.StaticContext) (*types.Players, error) {
    var players types.Players
    if err := ctx.Get(PLAYERS, &players); err != nil {return nil, err}
    return &players, nil
}

func (c *Huyilla) getActivePlayers (ctx contract.StaticContext) ([]*types.PlayerDetails, error) {
    players, err := c.getPlayers(ctx)
    if err != nil { return nil, err }

    var activePlayers []*types.PlayerDetails

    for _, player := range players.Players {
        if player.LoggedIn {
            entity, err := c.getEntity(ctx, player.EntityId)
            if err != nil { return nil, err }
            activePlayers = append(activePlayers, &types.PlayerDetails{Player: player, Entity: entity})
        }
    }

    return activePlayers, nil
}

//
// func (c *Huyilla) getPlayersDetails (ctx contract.StaticContext) ([]*types.PlayerDetails, error) {
//     players, err := c.getPlayers(ctx)
//     if err != nil { return nil, err }
//
//     details := make([]*types.PlayerDetails, len(players.Players))
//
//     for i, player := range players.Players {
//         entity, err := c.getEntity(ctx, player.Id)
//         if err != nil { return nil, err }
//         details[i] = &types.PlayerDetails{Player: player, Entity: entity}
//     }
//
//     return &players, nil
// }
