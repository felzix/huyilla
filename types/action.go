package types

import (
	"fmt"
	"github.com/pkg/errors"
)

type ActionType uint64

const (
	MoveActionType ActionType = iota + 1
	DigActionType
	PlaceActionType
	DropActionType
)

type ActionSpecific uint64

const (
	WhereTo ActionSpecific = iota + 1
	Tool
	GivenItem
	Target
)

type Action struct {
	Type   ActionType
	Issuer EntityId // player entity
	Actor  EntityId
	Specifics
}
type Specifics map[ActionSpecific]interface{}

type SerializableAction struct {
	Type      ActionType
	Issuer    EntityId // player entity
	Actor     EntityId
	Specifics SerializableSpecifics
}
type SerializableSpecifics map[ActionSpecific][]byte

func (a Action) Marshal() ([]byte, error) {
	specifics := make(SerializableSpecifics, len(a.Specifics))

	for key, value := range a.Specifics {
		switch key {
		case WhereTo:
			if blob, err := value.(AbsolutePoint).Marshal(); err != nil {
				return nil, err
			} else {
				specifics[key] = blob
			}
		default:
			return nil, errors.Errorf("Bad actions.specifics key %s", key)
		}
	}

	sa := SerializableAction{
		Type:      a.Type,
		Issuer:    a.Issuer,
		Actor:     a.Actor,
		Specifics: specifics,
	}
	return ToBytes(&sa)
}

func (a *Action) Unmarshal(input []byte) error {
	sa := SerializableAction{}
	if err := FromBytes(input, &sa); err != nil {
		return err
	}
	a.Actor = sa.Actor
	a.Issuer = sa.Issuer
	a.Type = sa.Type
	a.Specifics = make(Specifics, len(sa.Specifics))
	for key, blob := range sa.Specifics {
		switch key {
		case WhereTo:
			p := AbsolutePoint{}
			if err := p.Unmarshal(blob); err != nil {
				return err
			}
			a.Specifics[key] = p
		}
	}
	return nil
}

func NewMoveAction(player, actor EntityId, whereTo AbsolutePoint) *Action {
	return &Action{
		Type:   MoveActionType,
		Issuer: player,
		Actor:  actor,
		Specifics: Specifics{
			WhereTo: whereTo,
		},
	}
}

func (a Action) Apply(world *World) (bool, error) {
	entity, err := world.Entity(a.Issuer)
	if err != nil {
		return false, err
	} else if entity == nil {
		return false, errors.New(fmt.Sprintf(`Entity "%d" does not exist`, a.Issuer))
	}

	if len(entity.PlayerName) == 0 {
		return false, errors.New("entity does not have an associated player")
	}

	applier, err := ChooseActionApplier(a.Type)
	if err != nil {
		return false, nil
	}

	return applier(world, a)
}

type ActionApplier func(*World, Action) (bool, error)

func ChooseActionApplier(actionType ActionType) (ActionApplier, error) {
	switch actionType {
	case MoveActionType:
		return ApplyMoveAction, nil
	case DigActionType:
		return ApplyDigAction, nil
	case PlaceActionType:
		return ApplyPlaceAction, nil
	case DropActionType:
		return ApplyDropAction, nil
	default:
		return nil, errors.Errorf("Invalid action %d", uint64(actionType))
	}
}

func ApplyMoveAction(world *World, action Action) (bool, error) {
	actor, err := world.Entity(action.Actor)
	if err != nil {
		return false, err
	}

	if err := world.RemoveEntityFromChunk(actor.Id, actor.Location.Chunk); err != nil {
		return false, err
	}

	var whereTo AbsolutePoint
	switch s := action.Specifics[WhereTo].(type) {
	case AbsolutePoint:
		whereTo = s
	default:
		return false, errors.Errorf("Action property WhereTo is not an AbsolutePoint")
	}

	actor.Location = whereTo
	if err := world.SetEntity(actor.Id, actor); err != nil {
		return false, err
	}

	if err := world.AddEntityToChunk(actor); err != nil {
		return false, err
	}

	return true, nil
}

func ApplyDigAction(world *World, action Action) (bool, error) {
	return false, nil
}

func ApplyPlaceAction(world *World, action Action) (bool, error) {
	return false, nil
}

func ApplyDropAction(world *World, action Action) (bool, error) {
	return false, nil
}
