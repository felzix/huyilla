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

			viewDiameter := 3
			topbarHeight := 2
			boardSize := C.CHUNK_SIZE
			totalWidth := boardSize * viewDiameter
			totalheight := topbarHeight + (boardSize * viewDiameter)

			var child *react.Child
			if client.world.age == 0 {
				child = react.NewChild(react.Label(), "loading", maxWidth, maxHeight, react.Properties{
					"label": "Loading world from engine. Please wait.",
				})
			} else if totalheight > maxHeight || totalWidth > maxWidth {
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
									Key:     "",
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
									X:      0,
									Y:      0,
									Width:  maxWidth,
									Height: topbarHeight,
								},
								{
									Element: Tiles(),
									Key:     "",
									Props: react.Properties{
										"client": client,
										"center": client.player.Entity.Location,
									},
									X:      0,
									Y:      topbarHeight,
									Width:  boardSize * viewDiameter,
									Height: boardSize * viewDiameter,
								},
							},
						}, nil
					},
					HandleKeyFn: func(element *react.ReactElement, e *tcell.EventKey) (bool, error) {
						var to *types.AbsolutePoint

						switch e.Rune() {
						case 'w': // move up
							to = client.player.Entity.Location.Derive(0, -1, 0, C.CHUNK_SIZE)
						case 's': // move down
							to = client.player.Entity.Location.Derive(0, +1, 0, C.CHUNK_SIZE)
						case 'a': // move left
							to = client.player.Entity.Location.Derive(-1, 0, 0, C.CHUNK_SIZE)
						case 'd': // move right
							to = client.player.Entity.Location.Derive(+1, 0, 0, C.CHUNK_SIZE)
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

func drawMissingChunk(result react.DrawResult, localX, localY, width, height int) {
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			drawMissingTile(result, x, y, localX, localY)
		}
	}
}

func drawMissingTile(result react.DrawResult, x, y, localX, localY int) {
	result.Region.Cells[x+localX][y+localY] = react.Cell{
		R:     ' ',
		Style: tcell.StyleDefault.Background(tcell.ColorDarkGray),
	}
}

func drawChunk(result react.DrawResult, localX, localY, width, height, zLevel int, chunk *types.DetailedChunk) {
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			drawTile(result, x, y, localX, localY, zLevel, chunk)
		}
	}
}

func drawTile(result react.DrawResult, x, y, localX, localY, zLevel int, chunk *types.DetailedChunk) {
	index := (x * C.CHUNK_SIZE * C.CHUNK_SIZE) + (y * C.CHUNK_SIZE) + zLevel
	voxel := types.Voxel(chunk.Voxels[index])
	rune_ := voxelToRune(voxel)
	result.Region.Cells[x+localX][y+localY] = react.Cell{
		R:     rune_,
		Style: tcell.StyleDefault,
	}
}

func drawEntitiesForChunk(result react.DrawResult, localX, localY, width, height, zLevel int, chunk *types.DetailedChunk) {
	for _, entity := range chunk.Entities {
		x := int(entity.Location.Voxel.X)
		y := int(entity.Location.Voxel.Y)
		z := int(entity.Location.Voxel.Z)

		if x >= width || y >= height || z != zLevel {
			continue
		}

		drawEntity(result, x, y, localX, localY, entity)
	}
}

func drawEntity(result react.DrawResult, x, y, localX, localY int, entity *types.Entity) {
	result.Region.Cells[x+localX][y+localY] = react.Cell{
		R:     entityToRune(entity),
		Style: tcell.StyleDefault,
	}
}

func Tiles() *react.ReactElement {
	return &react.ReactElement{
		Type: "Tiles",
		DrawFn: func(r *react.ReactElement, maxWidth, maxHeight int) (*react.DrawResult, error) {
			client := r.Props["client"].(*Client)
			center := r.Props["center"].(*types.AbsolutePoint)

			result := react.DrawResult{
				Region: react.NewRegion(0, 0, maxWidth, maxHeight),
			}

			localX := 0
			localY := 0

			for chunkX := -1; chunkX < 2; chunkX++ {
				width := C.CHUNK_SIZE
				if width > maxWidth-localX {
					width = maxWidth - localX
				}

				for chunkY := -1; chunkY < 2; chunkY++ {
					height := C.CHUNK_SIZE
					if height > maxHeight-localY {
						height = maxHeight - localY
					}

					point := center.Derive(int64(chunkX*C.CHUNK_SIZE), int64(chunkY*C.CHUNK_SIZE), 0, C.CHUNK_SIZE)
					chunk := client.world.chunks[*types.NewComparablePoint(point.Chunk)]

					zLevel := int(point.Voxel.Z)

					if chunk == nil {
						Log.Warningf("react:Tiles: chunk missing @ %v", point.Chunk)
						drawMissingChunk(result, localX, localY, width, height)
					} else {
						drawChunk(result, localX, localY, width, height, zLevel, chunk)
						drawEntitiesForChunk(result, localX, localY, width, height, zLevel, chunk)
					}

					localY += height
				}

				localY = 0
				localX += width
			}
			return &result, nil
		},
	}
}
