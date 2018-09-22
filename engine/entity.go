package engine


type Entity struct {
    eType EntityType
    properties map[string]interface{}
}

type EntityType uint

func MakeEntity (eType EntityType) *Entity {
    return &Entity{eType: eType}
}

func (e *Entity) GetProperty (content *Content, prop string) interface{} {
    value := e.properties[prop]
    if value == nil {
        value = content.EP[e.eType]
    }
    return value
}
