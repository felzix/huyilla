package main

import (
    "encoding/json"
    "github.com/pkg/errors"
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

func (c *Huyilla) GetPlayer (ctx contract.StaticContext, req *types.PlayerName) (*types.PlayerDetails, error) {
    player, err := c.getPlayer(ctx, req.Name)
    if err != nil {return nil, errors.Wrap(err, "player not found")}

    entity, err := c.getEntity(ctx, player.Id)
    if err != nil {return nil, errors.Wrap(err, "entity not found")}

    details := &types.PlayerDetails{Player: player, Entity: entity}
    return details, nil
}

func (c *Huyilla) CreatePlayer(ctx contract.Context, req *types.PlayerName) error {
    players, err := c.getPlayers(ctx)
    if err != nil {return err}

    player := players.Players[req.Name]
    if player != nil {
        return errors.New(fmt.Sprintf(`Name "%v" is taken. Try another.`, req.Name))
    }

    player = &types.Player{
        Name:    req.Name,
        Address: []byte(ctx.Message().Sender.Local),
    }
    players.Players[player.Name] = player

    err = ctx.Set(PLAYERS, players)
    if err != nil {return err}

    ctx.Logger().Info("Created player", "name", player.Name, "address", player.Address)

    emitMsg := struct {
        Method string
        Owner  string
        Addr   []byte
    }{"CreatePlayer", player.Name, player.Address}
    emitMsgJSON, err := json.Marshal(emitMsg)
    if err != nil {return err}

    ctx.EmitTopics(emitMsgJSON, "huyilla:" + emitMsg.Method)
    return nil
}

func (c *Huyilla) getPlayer (ctx contract.StaticContext, name string) (*types.Player, error) {
    var players types.Players
    if err := ctx.Get(PLAYERS, &players); err != nil {return nil, err}
    player := players.Players[name]
    return player, nil
}

func (c *Huyilla) getPlayers (ctx contract.StaticContext) (*types.Players, error) {
    var players types.Players
    if err := ctx.Get(PLAYERS, &players); err != nil {return nil, err}
    return &players, nil
}
