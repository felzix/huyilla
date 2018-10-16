package main

import (
    "github.com/felzix/huyilla/types"
    "github.com/loomnetwork/go-loom/plugin"
    contract "github.com/loomnetwork/go-loom/plugin/contractpb"
)


func main() {
    plugin.Serve(Contract)
}

type Huyilla struct {}


var Contract = contract.MakePluginContract(&Huyilla{})


func (c *Huyilla) Meta () (plugin.Meta, error) {
    return plugin.Meta{
        Name:    "Huyilla",
        Version: "0.0.1",
    }, nil
}

var AGE = []byte("Age")
var CONFIG = []byte("Config")
var PLAYERS = []byte("Players")

func (c *Huyilla) Init (ctx contract.Context, req *plugin.Request) error {
    ctx.Set(AGE, &types.Age{Ticks: 1})  // starts at 1 because 0 counts as non-existent

    config := &types.Config{
        Options: &types.PrimitiveMap{
            Map: map[string]*types.Primitive{
                "PlayerCap": &types.Primitive{Value: &types.Primitive_Int{Int: 10}},
            },
        },
    }
    ctx.Set(CONFIG, config)

    return nil
}

func (c *Huyilla) GetAge (ctx contract.StaticContext, req *plugin.Request) (*types.Age, error) {
    var age = &types.Age{}
    err := ctx.Get(AGE, age)
    return age, err
}

func (c *Huyilla) incrementAge (ctx contract.Context, req *plugin.Request) (*types.Age, error) {
    age, err := c.GetAge(ctx, req)
    if err != nil { return nil, err }

    age.Ticks ++
    err = ctx.Set(AGE, age)
    return age, err
}

func (c *Huyilla) GetConfig (ctx contract.StaticContext, req *plugin.Request) (*types.Config, error) {
    return c.getConfig(ctx)
}

func (c *Huyilla) SetConfigOptions (ctx contract.Context, req *types.PrimitiveMap) error {
    config, err := c.getConfig(ctx)
    if err != nil { return err }

    for k,v := range req.Map {
        config.Options.Map[k] = v
    }

    return ctx.Set(CONFIG, config)
}


func (c *Huyilla) getConfig (ctx contract.StaticContext) (*types.Config, error) {
    var config = &types.Config{}
    err := ctx.Get(CONFIG, config)
    return config, err
}
