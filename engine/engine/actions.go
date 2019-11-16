package engine

import (
	"github.com/felzix/huyilla/types"
)

func (engine *Engine) RegisterAction(action *types.Action) {
	engine.Lock()
	defer engine.Unlock()

	engine.Actions = append(engine.Actions, action)
}
