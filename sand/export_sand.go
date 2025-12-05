package sand

import (
	// "fmt"
	"time"

	"timer/timer/types"

	"github.com/rivo/tview"
)

const (
	width  = 5
	height = 50
	blocks = 250

// // 5 * 50 == 250
// // 2 min == 120 secs
// // 250 blocks / 120 secs = 2.08 ~~ 2 blocks per sec
)

const (
	Empty types.Cell = iota
	Sand
)

func NewSandGrid(app *tview.Application, timers *[]*types.Timer) *tview.Grid {

	grid := tview.NewGrid().SetRows(0).SetColumns(0, 0)
	for i, t := range *timers {
		t.V = tview.NewTextView().SetDynamicColors(true)
		grid.AddItem(t.V, 0, i, 1, 1, 0, 0, false)

		// initialize grid
		t.G = make([][]types.Cell, height)
		for i := range t.G {
			t.G[i] = make([]types.Cell, width)
		}

		// Function to render gridRight as ANSI-colored string
		t.RF = func() string {
			s := ""
			for y := 0; y < height; y++ {
				for x := 0; x < width; x++ {
					if t.G[y][x] == Sand {
						s += "[blue]█[-]"
					} else {
						s += " "
					}
				}
				s += "\n"
			}
			return s
		}

		// animation loop
		t.PH = func(grid [][]types.Cell, render func() string) {
			// seconds / block (seconds per block)
			s := t.L.Seconds()
			mpt := (s / blocks) * 1000 // * 1000 to get Millisecond
			tick := time.NewTicker(time.Duration(mpt) * time.Millisecond)
			x := 0
			for {
				// Update sand physics
				for y := height - 2; y >= 0; y-- {
					for x := 0; x < width; x++ {
						if grid[y][x] == Sand {
							if grid[y+1][x] == Empty {
								grid[y+1][x], grid[y][x] = Sand, Empty
							} else {
								if x > 0 && grid[y+1][x-1] == Empty {
									grid[y+1][x-1], grid[y][x] = Sand, Empty
								} else if x < width-1 && grid[y+1][x+1] == Empty {
									grid[y+1][x+1], grid[y][x] = Sand, Empty
								}
							}
						}
					}
				}
				select {
				case <-tick.C:
					grid[0][x%width] = Sand
					x++
				default:
				}

				app.QueueUpdateDraw(func() {
					t.V.SetText(render())
				})

				if mpt > 100 {
					time.Sleep(100 * time.Millisecond)
				}
			}
		}
	}

	return grid
}

// func Render() string {
// 	// Function to render gridRight as ANSI-colored string
// 	render := func() string {
// 		s := ""
// 		for y := 0; y < height; y++ {
// 			for x := 0; x < width; x++ {
// 				if gridRight[y][x] == Sand {
// 					s += "[blue]█[-]"
// 				} else {
// 					s += " "
// 				}
// 			}
// 			s += "\n"
// 		}
// 		return s
// 	}
//
// 	mpt := (10.0 / 250.0) * 1000
// 	tick := time.NewTicker(time.Duration(mpt) * time.Millisecond)
// 	x := 0
// 	// Animation loop
// 	go func() {
// 		for {
// 			// Update sand physics
// 			for y := height - 2; y >= 0; y-- {
// 				for x := 0; x < width; x++ {
// 					if gridRight[y][x] == Sand {
// 						if gridRight[y+1][x] == Empty {
// 							gridRight[y+1][x], gridRight[y][x] = Sand, Empty
// 						} else {
// 							if x > 0 && gridRight[y+1][x-1] == Empty {
// 								gridRight[y+1][x-1], gridRight[y][x] = Sand, Empty
// 							} else if x < width-1 && gridRight[y+1][x+1] == Empty {
// 								gridRight[y+1][x+1], gridRight[y][x] = Sand, Empty
// 							}
// 						}
// 					}
// 				}
// 			}
//
// 			// Occasionally add new sand at the top
// 			// for x := 0; x < width; x++ {
// 			select {
// 			case <-tick.C:
// 				gridRight[0][x%width] = Sand
// 				x++
// 			default:
// 			}
// 			// }
// 			// for x := 0; x < width; x++ {
// 			// 	if rand.Float32() < 0.05 {
// 			// 		gridRight[0][x] = Sand
// 			// 	}
// 			// }
//
// 			app.QueueUpdateDraw(func() {
// 				rightTextView.SetText(render())
// 			})
//
// 			if mpt > 100 {
// 				time.Sleep(100 * time.Millisecond)
// 			}
// 		}
// 	}()
// 	}
