package types

type PlayerDetails struct {
	Player *Player
	Entity *Entity
}

func NewPlayerDetails(player *Player, entity *Entity) *PlayerDetails {
	return &PlayerDetails{
		Player: player,
		Entity: entity,
	}
}

func (p PlayerDetails) Marshal() ([]byte, error) {
	return ToBytes(p)
}

func (p *PlayerDetails) Unmarshal(blob []byte) error {
	return FromBytes(blob, &p)
}
