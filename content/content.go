package content

import "github.com/felzix/huyilla/types"

var ContentDefinitions = types.Content{
	E: EntityDefinitions,
	F: FormDefinitions,
	M: MaterialDefinitions,
}

var ENTITY = make(map[string]uint64)
var FORM = make(map[string]uint64)
var MATERIAL = make(map[string]uint64)

func PopulateContentNameMaps() {
	for id, def := range ContentDefinitions.E {
		ENTITY[def.Name] = id
	}
	for id, def := range ContentDefinitions.F {
		FORM[def.Name] = id
	}
	for id, def := range ContentDefinitions.M {
		MATERIAL[def.Name] = id
	}
}
