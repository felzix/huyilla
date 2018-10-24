package main

import (
    "github.com/felzix/huyilla/types"
    "github.com/loomnetwork/go-loom/plugin"
    contract "github.com/loomnetwork/go-loom/plugin/contractpb"
)


var AGE = []byte("Age")


func (c *Huyilla) GetAge (ctx contract.StaticContext, req *plugin.Request) (*types.Age, error) {
    return c.getAge(ctx)
}

func (c *Huyilla) incrementAge (ctx contract.Context) (*types.Age, error) {
    age, err := c.getAge(ctx)
    if err != nil { return nil, err }

    age.Ticks ++
    err = ctx.Set(AGE, age)
    return age, err
}

func (c *Huyilla) getAge (ctx contract.StaticContext) (*types.Age, error) {
    var age = &types.Age{}
    err := ctx.Get(AGE, age)
    return age, err
}