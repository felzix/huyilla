package engine


type Action interface {
	Type() ActionType
	Apply(*World)
}
type ActionType uint8

type Move struct {
    Who       UniqueId
    WhereFrom Point
    WhereTo   Point
    WhereAt   At
}

func (act Move) Type () ActionType {
    return 1
}

// TODO validation
func (act Move) Apply (world *World) bool {
    entity := world.GetEntity(act.Who, act.WhereFrom)

    if entity == nil {
        return false
    }

    world.RemoveEntity(act.Who, act.WhereFrom)
    world.SetEntity(entity, act.WhereTo, act.WhereAt)

    return true
}
