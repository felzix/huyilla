package main

import (
	"fmt"
	react "github.com/felzix/go-curses-react"
	C "github.com/felzix/huyilla/constants"
	"github.com/felzix/huyilla/types"
	"github.com/gdamore/tcell"
)


const (
	VIEWMODE_INTRO = iota
	VIEWMODE_GAME
)

func MakeApp() *react.ReactElement {
	root := &react.ReactElement{
		State: react.State{
			"mode": VIEWMODE_INTRO,
		},
		DrawFn: func(r *react.ReactElement, maxWidth, maxHeight int) (*react.DrawResult, error) {
			client := r.Props["client"].(*Client)
			mode := r.State["mode"].(int)

			var element *react.ReactElement
			var props react.Properties
			switch mode {
			case VIEWMODE_INTRO:
				element = Intro()
				props = react.Properties{
					"client": client,
					"nextMode": func() {
						r.State["mode"] = VIEWMODE_GAME
					},
				}
			case VIEWMODE_GAME:
				element = GameBoard()
				props = react.Properties{
					"client": client,
				}
			}

			result := react.DrawResult{
				Elements: []react.Child{
					*react.NewChild(element, string(mode), maxWidth, maxHeight, props),
				}}
			return &result, nil
		},
	}

	return root
}

func Intro() *react.ReactElement {
	return &react.ReactElement{
		Type: "Intro",
		DrawFn: func(r *react.ReactElement, maxWidth, maxHeight int) (*react.DrawResult, error) {
			client := r.Props["client"].(*Client)
			nextMode := r.Props["nextMode"].(func())

			child := react.NewChild(react.HorizontalLayout(), "", maxWidth, maxHeight, react.Properties{
				"children": []*react.Child{
					react.ManagedChild(react.Label(), "hello", react.Properties{
						"label": "Hello!",
					}),
					react.ManagedChild(react.Label(), "blank", react.Properties{
						"label": "",
					}),
					react.ManagedChild(react.TextEntry(), "", react.Properties{
						"label": "Enter username",
						"whenFinished": func(username string) error {
							client.username = username

							if err := client.Auth(); err != nil {
								return err
							}

							nextMode()
							return nil
						},
						// TODO when TextEntry can do validation, reject empty or taken username
					}),
				},
			})
			result := react.DrawResult{
				Elements: []react.Child{*child},
			}
			return &result, nil
		},
	}
}

func GameBoard() *react.ReactElement {
	return &react.ReactElement{
		Type: "GameBoard",
		DrawFn: func(r *react.ReactElement, maxWidth, maxHeight int) (*react.DrawResult, error) {
			client := r.Props["client"].(*Client)

			topbarHeight := 2
			boardSize := C.CHUNK_SIZE
			totalheight := topbarHeight + boardSize

			var child *react.Child
			if client.world.age == 0 {
				child = react.NewChild(react.Label(), "loading", maxWidth, maxHeight, react.Properties{
					"label": "Loading world from engine. Please wait.",
				})
			} else if totalheight > maxHeight || boardSize > maxWidth {
				child = react.NewChild(react.Label(), "screen-too-small", maxWidth, maxHeight, react.Properties{
					"label": "Terminal screen too small",
				})
			} else {
				container := &react.ReactElement{
					Type: "gameboard-container",
					Key:  "only",
					DrawFn: func(element *react.ReactElement, maxWidth int, maxHeight int) (*react.DrawResult, error) {
						return &react.DrawResult{
							Elements: []react.Child{
								{
									Element: react.HorizontalLayout(),
									Key: "",
									Props: react.Properties{
										"children": []*react.Child{
											react.ManagedChild(react.Label(), "debug-bar", react.Properties{
												"label": fmt.Sprintf("%d", client.world.age),
											}),
											react.ManagedChild(react.Label(), "blank", react.Properties{
												"label": "",
											}),
										},
									},
									X: 0,
									Y: 0,
									Width: maxWidth,
									Height: 2,
								},
								{
									Element: Tiles(),
									Key: "",
									Props: react.Properties{
										"client":   client,
										"absPoint": client.player.Entity.Location,
									},
									X: 0,
									Y: 2,
									Width: boardSize,
									Height: boardSize,
								},
							},
						}, nil
					},
					HandleKeyFn: func(element *react.ReactElement, e *tcell.EventKey) (bool, error) {
						var to *types.AbsolutePoint

						switch e.Rune() {
						case 'w': // move up
							to = client.player.Entity.Location
							to.Voxel.Y--
						case 's': // move down
							to = client.player.Entity.Location
							to.Voxel.Y++
						case 'a': // move left
							to = client.player.Entity.Location
							to.Voxel.X--
						case 'd': // move right
							to = client.player.Entity.Location
							to.Voxel.X++
						}

						if to != nil {
							err := client.api.IssueMoveAction(to)
							return false, err
						} else {
							return true, nil
						}
					},
				}

				child = react.NewChild(container, "gameboard", maxWidth, maxHeight, nil)
			}

			result := react.DrawResult{
				Elements: []react.Child{*child},
			}
			return &result, nil
		},
	}
}

func Tiles() *react.ReactElement {
	return &react.ReactElement{
		Type: "Tiles",
		DrawFn: func(r *react.ReactElement, maxWidth, maxHeight int) (*react.DrawResult, error) {
			client := r.Props["client"].(*Client)
			absPoint := r.Props["absPoint"].(*types.AbsolutePoint)

			chunk := client.world.chunks[*types.NewComparablePoint(absPoint.Chunk)]
			zLevel := absPoint.Voxel.Z

			width := C.CHUNK_SIZE
			if width > maxWidth {
				width = maxWidth
			}
			height := C.CHUNK_SIZE
			if height > maxHeight {
				height = maxHeight
			}

			result := react.DrawResult{
				Region: react.NewRegion(0, 0, maxWidth, maxHeight),
			}

			if chunk == nil {
				for y := 0; y < height; y++ {
					for x := 0; x < width; x++ {
						result.Region.Cells[x][y] = react.Cell{
							R:     ' ',
							Style: tcell.StyleDefault.Background(tcell.ColorDarkGray),
						}
					}
				}
			} else {
				for y := 0; y < height; y++ {
					for x := 0; x < width; x++ {
						index := (x * C.CHUNK_SIZE * C.CHUNK_SIZE) + (y * C.CHUNK_SIZE) + int(zLevel)
						ch := voxelToRune(chunk.Voxels[index])
						result.Region.Cells[x][y] = react.Cell{
							R:     ch,
							Style: tcell.StyleDefault,
						}
					}
				}

				for _, entity := range chunk.Entities {
					x := entity.Location.Voxel.X
					y := entity.Location.Voxel.Y
					z := entity.Location.Voxel.Z

					if z == zLevel {
						result.Region.Cells[x][y] = react.Cell{
							R: entityToRune(entity),
							Style: tcell.StyleDefault,
						}
					}

				}
			}
			return &result, nil
		},
	}
}
