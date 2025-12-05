package types

import (
	"time"

	"github.com/rivo/tview"
)

type Cell int

type Timer struct {
	N string        //name
	L time.Duration // timer length // time.Duration to avoid casting
	S time.Time
	T *time.Timer // the Timer for the custom type (timer)
	A bool        // active timer bool // true == active
	P time.Time   // timer pause marker

	V  *tview.TextView
	G  [][]Cell
	RF func() string                 // render func
	PH func([][]Cell, func() string) // physics func
}
