package types

func NewMove(player string, whereTo *AbsolutePoint) *Action{
	return &Action{
		PlayerName: player,
		Action: &Action_Move{
			Move: &Action_MoveAction{
				WhereTo: whereTo,
			},
		},
	}
}

