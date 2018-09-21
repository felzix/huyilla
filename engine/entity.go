package engine


type Entity struct {
    eType EntityType
    properties map[string]interface{}
}

type EntityType uint
type EntityProperty uint

func (e *Entity) GetProperty (prop string) interface{} {
    value := e.properties[prop]
    if value == nil {
        value = 12
    }
    return value
}
