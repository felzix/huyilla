package main

import (
	"github.com/felzix/huyilla/types"
	"github.com/loomnetwork/go-loom/plugin"
	contract "github.com/loomnetwork/go-loom/plugin/contractpb"
)

var CONFIG = []byte("Config")

func (c *Huyilla) GetConfig(ctx contract.StaticContext, req *plugin.Request) (*types.Config, error) {
	return c.getConfig(ctx)
}

func (c *Huyilla) SetConfigOptions(ctx contract.Context, req *types.PrimitiveMap) error {
	config, err := c.getConfig(ctx)
	if err != nil {
		return err
	}

	for k, v := range req.Map {
		config.Options.Map[k] = v
	}

	return ctx.Set(CONFIG, config)
}

func (c *Huyilla) getConfig(ctx contract.StaticContext) (*types.Config, error) {
	var config = &types.Config{}
	err := ctx.Get(CONFIG, config)
	return config, err
}
