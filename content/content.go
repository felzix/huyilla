package content

import "github.com/felzix/huyilla/types"

var Definitions = types.Content{
	E: EntityDefinitions,
	F: FormDefinitions,
	M: MaterialDefinitions,
}

var ENTITY = make(map[string]types.EntityType)
var FORM = make(map[string]types.Form)
var MATERIAL = make(map[string]types.Material)

func init() {
	for id, def := range Definitions.E {
		ENTITY[def.Name] = id
	}
	for id, def := range Definitions.F {
		FORM[def.Name] = id
	}
	for id, def := range Definitions.M {
		MATERIAL[def.Name] = id
	}
}
