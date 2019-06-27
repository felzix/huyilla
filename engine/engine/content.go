package engine

import (
	"github.com/felzix/huyilla/content"
	"github.com/felzix/huyilla/types"
)

var ENTITY = content.ENTITY
var FORM = content.FORM
var MATERIAL = content.MATERIAL

func (engine *Engine) GetContent() *types.Content {
	return &content.ContentDefinitions
}
