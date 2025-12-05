package main

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"timer/timer/sand"
	"timer/timer/types"
)

var timers []*types.Timer
var currentView string

// var toastNotificationLen = 2 * time.Second

// type timer struct {
// 	n string        //name
// 	l time.Duration // timer length // time.Duration to avoid casting
// 	s time.Time
// 	t *time.Timer // the Timer for the custom type (timer)
// 	a bool        // active timer bool // true == active
// 	p time.Time   // timer pause marker
// }

func getTimer(timerName string) (bool, *types.Timer) {
	for _, t := range timers {
		if t.N == timerName {
			return true, t
		}
	}
	return false, nil
}

func showTimeRemaining(timerName string, toast func(string, time.Duration)) {
	_, t := getTimer(timerName)
	if !t.A && !t.P.IsZero() {
		toast(fmt.Sprintf("timer \"%v\" is paused with %v remaining", t.N, t.L.Round(time.Second)), 2*time.Second)
	} else if t.A {
		toast(fmt.Sprintf("timer \"%v\" has %v remaining\n", t.N, time.Until(t.S.Add(t.L)).Round(time.Second)), 2*time.Second)
	} else {
		toast(fmt.Sprintf("timer \"%v\" is not active\n", t.N), 2*time.Second)
	}
}

func resumeTimer(timerName string, toast func(string, time.Duration), timerFinishFunc func(string)) {
	_, t := getTimer(timerName)
	if !t.A && !t.P.IsZero() {
		t.S = time.Now()                   // new timer start time
		t.T = time.AfterFunc(t.L, func() { // new Afterfunc
			timerFinishFunc(timerName)
			t.A = false // time not active after it finishes
		})
		t.A = true        // timer active again after resuming
		t.P = time.Time{} // clear paused marker for printing purposes
		toast(fmt.Sprintf("timer \"%v\" resumed", t.N), 2*time.Second)
	} else {
		toast(fmt.Sprintf("timer \"%v\" is not paused", t.N), 2*time.Second)
	}
}

func stopTimer(timerName string, toast func(string, time.Duration)) {
	_, t := getTimer(timerName)
	if !t.T.Stop() && t.P.IsZero() {
		toast(fmt.Sprintf("timer \"%v\" has already fired or been stopped", t.N), 2*time.Second)
	} else {
		toast(fmt.Sprintf("timer \"%v\" stopped", t.N), 2*time.Second)
		t.P = time.Time{} // clear paused marker in case paused timer gets stopped
		t.A = false
	}
}

func pauseTimer(timerName string, toast func(string, time.Duration)) {
	_, t := getTimer(timerName)
	if !t.T.Stop() { // timer cannot be stopped for 1 or 2 reasons
		if t.P.IsZero() { // either a timer has no pause marker
			toast(fmt.Sprintf("timer \"%v\" has already fired or been stopped", t.N), 2*time.Second)
		} else { // or timer has pause marker
			toast(fmt.Sprintf("timer \"%v\" is already paused", t.N), 2*time.Second)
		}
	} else {
		t.P = time.Now()                 // mark when timer was paused
		timerLen := t.L - (t.P.Sub(t.S)) // find new timer length based on time elapsed - total time
		t.L = timerLen
		t.A = false // timer not active for printing purposes
		toast(fmt.Sprintf("timer \"%v\" paused", t.N), 2*time.Second)
	}
}

func createTimer(timerName string, length int, interval time.Duration, timerFinishFunc func(string)) {
	timerLength := time.Duration(length) * interval
	t := types.Timer{
		N: timerName,
		L: timerLength,
		S: time.Now(),
		A: true,
	}
	t.T = time.AfterFunc(timerLength, func() {
		timerFinishFunc(timerName)
		t.A = false
	})
	timers = append(timers, &t)
}

func main() {
	newPrimitive := func(text string) tview.Primitive {
		return tview.NewTextView().
			SetTextAlign(tview.AlignCenter).
			SetText(text)
	}

	app := tview.NewApplication()

	inputField := tview.NewInputField().
		SetLabel("timer name: ").
		SetFieldWidth(0).
		SetFieldTextColor(tcell.ColorBlack).
		SetFieldBackgroundColor(tcell.ColorWhite)

	header := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(newPrimitive("Timers"), 2, 0, false).
		AddItem(inputField, 0, 1, true)

	// different views for middle row
	//-----------------------------------------------------------------------------------------------------------------------------------
	// timer list view
	timerList := tview.NewList().
		ShowSecondaryText(false).
		SetHighlightFullLine(false).SetSelectedBackgroundColor(tcell.ColorBeige)

	updateList := func(filter string) {
		timerList.Clear()
		for _, timer := range timers {
			if strings.Contains(strings.ToLower(timer.N), strings.ToLower(filter)) {
				timerList.AddItem(fmt.Sprintf("[green]%v[-]", timer.N), "", 0, nil)
			}
		}
	}

	// timer creation view
	timerNameInput := tview.NewInputField().
		SetLabel("new timer name: ").
		SetFieldWidth(10).
		SetFieldTextColor(tcell.ColorBlack).
		SetFieldBackgroundColor(tcell.ColorWhite)

	timerDurationInput := tview.NewInputField().
		SetLabel("new timer length: ").
		SetFieldWidth(10).
		SetFieldTextColor(tcell.ColorBlack).
		SetFieldBackgroundColor(tcell.ColorWhite)

	timerIntervalInput := tview.NewList().
		ShowSecondaryText(false).
		SetHighlightFullLine(false)

	timerIntervals := []string{"millisecond", "second", "minute", "hour"}
	for _, interval := range timerIntervals {
		timerIntervalInput.AddItem(interval, "", 0, nil)
	}

	createTimerView := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(timerNameInput, 0, 1, false).
		AddItem(timerDurationInput, 0, 1, false).
		AddItem(timerIntervalInput, 0, 1, false)
	//-----------------------------------------------------------------------------------------------------------------------------------

	footer := tview.NewFlex()
	searchCmds := []string{"[enter] list/create t", "[esc] sand"}
	listCmds := []string{"create t", "stop t", "[enter] options", "[esc] search"}
	timerCmds := []string{"pause t", "resume r", "show t", "[esc] list"}

	showSearchCmds := func() {
		footer.Clear()
		for _, c := range searchCmds {
			footer.AddItem(newPrimitive(c), 0, 1, false)
		}
	}

	showListCmds := func() {
		footer.Clear()
		for i, c := range listCmds {
			if i == 2 || i == 3 {
				footer.AddItem(newPrimitive(c), 0, 1, false)
			} else {
				footer.AddItem(newPrimitive(fmt.Sprintf("[%v] %v", string(c[0]), c)), 0, 1, false)
			}
		}
	}

	showTimerCmds := func() {
		footer.Clear()
		for i, c := range timerCmds {
			if i == 3 {
				footer.AddItem(newPrimitive(c), 0, 1, false)
			} else {
				footer.AddItem(newPrimitive(fmt.Sprintf("[%v] %v", string(c[0]), c)), 0, 1, false)
			}
		}
	}

	// grid to hold items (not included modals)
	grid := tview.NewGrid().
		SetRows(3, 0, 1).
		SetBorders(true).
		AddItem(header, 0, 0, 1, 1, 0, 0, false).
		AddItem(timerList, 1, 0, 1, 1, 0, 0, false).
		AddItem(footer, 2, 0, 1, 1, 0, 0, false)

	rootPage := tview.NewPages()
	rootPage.AddPage("mainPage", grid, true, true)
	rootPage.AddPage("createTimerViewPage", createTimerView, true, false)
	rootPage.AddPage("sandPage", sand.NewSandGrid(app, &timers), true, false)

	updateSand := func() {
		rootPage.RemovePage("sandPage")
		rootPage.AddPage("sandPage", sand.NewSandGrid(app, &timers), true, false)
	}

	currentView = "inputFieldView"

	updateGuidance := func() {
		switch currentView {
		case "inputFieldView":
			showSearchCmds()
		case "listView":
			showListCmds()
		case "timerOptionsView":
			showTimerCmds()
		}
	}
	updateGuidance()

	showToast := func(msg string, duration time.Duration) {
		toast := tview.NewTextView().SetText(msg).SetTextAlign(tview.AlignCenter)
		toast.SetBackgroundColor(tcell.ColorFireBrick)

		toastBox := tview.NewFlex().AddItem(nil, 0, 5, false).AddItem(tview.NewFlex().SetDirection(tview.FlexRow).AddItem(nil, 5, 0, false).AddItem(toast, 2, 0, false), 0, 1, false)

		// TODO test this current/setCurrent implementation
		// annoying to have inputField always retain focus after toast pages hides

		rootPage.AddPage("toastPage", toastBox, true, true)

		go func() {
			time.Sleep(duration)
			app.QueueUpdateDraw(func() {
				rootPage.HidePage("toastPage")
				if currentView != "inputFieldView" {
					app.SetFocus(timerList)
				} else {
					app.SetFocus(inputField)
				}
			})
			updateGuidance()
		}()
	}

	// -------------------------------------------------------------------------------------
	// TODO please refactor this sometime PLS!
	// i think need to make global string that is focus and then set that each time
	// then call a func at the bottom to change focus
	// the problem is trying to call focus on object that doesnt exist yet because
	// i want object created top to bottom similar to app
	// ex. header cannot set focus on footer because footer is created after header

	timerList.SetDoneFunc(func() {
		if currentView == "timerOptionsView" {
			currentView = "listView"
			showListCmds()
		} else {
			app.SetFocus(inputField)
			showSearchCmds()
		}
	})

	timerList.SetSelectedFunc(func(_ int, _ string, _ string, _ rune) {
		showTimerCmds()
		currentView = "timerOptionsView"
	})

	inputField.SetChangedFunc(func(text string) {
		updateList(text)
	})

	inputField.SetDoneFunc(func(key tcell.Key) {
		if key == tcell.KeyEnter {
			app.SetFocus(timerList)
			showListCmds()
		}
	})

	// -------------------------------------------------------------------------------------

	showList := func() {
		rootPage.SwitchToPage("mainPage")
		app.SetFocus(inputField)
		showSearchCmds()
		updateList("")
		currentView = "listView"
	}

	showDetail := func() {
		rootPage.SwitchToPage("createTimerViewPage")
		app.SetFocus(timerNameInput)
		currentView = "createTimerView"
	}

	// timer complete modal
	close := []string{"close", "hi"}
	timerComplete := tview.NewModal().
		AddButtons(close).
		SetText("Timer ____ is done!").
		SetDoneFunc(func(int, string) {
			// app.SetRoot(grid, false).SetFocus(inputField)
			app.SetRoot(rootPage, false).SetFocus(inputField)
			updateList("")
		})

	showTimerComplete := func(timerName string) {
		timerComplete.SetText(fmt.Sprintf("Timer \"%v\" is done!", timerName))
		app.SetRoot(timerComplete, true).Draw()
	}

	timerNameInput.SetDoneFunc(func(key tcell.Key) {
		if key == tcell.KeyTab {
			tn := timerNameInput.GetText()
			if tn != "" {
				app.SetFocus(timerDurationInput)
			}
		}
	})

	timerDurationInput.SetDoneFunc(func(key tcell.Key) {
		if key == tcell.KeyTab {
			_, err := strconv.Atoi(timerDurationInput.GetText())
			if err == nil {
				app.SetFocus(timerIntervalInput)
			}
		}
	})

	timerIntervalInput.SetSelectedFunc(func(ti int, _ string, _ string, _ rune) {
		// get data from input fields and then reset to default
		tn := timerNameInput.GetText()
		timerNameInput.SetText("")
		tl, _ := strconv.Atoi(timerDurationInput.GetText())
		timerDurationInput.SetText("")
		timerIntervalInput.SetCurrentItem(0)

		switch ti {
		case 0:
			createTimer(tn, tl, time.Millisecond, showTimerComplete)
		case 1:
			createTimer(tn, tl, time.Second, showTimerComplete)
		case 2:
			createTimer(tn, tl, time.Minute, showTimerComplete)
		case 3:
			createTimer(tn, tl, time.Hour, showTimerComplete)
		default:
			createTimer(tn, tl, time.Second, showTimerComplete)
		}
		showList()
		showToast(fmt.Sprintf("timer \"%v\" created", tn), 2*time.Second)
		updateSand()
	})

	app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyEnter {
			if timerIntervalInput.HasFocus() {
				return event
			} else {
				if len(timers) == 0 {
					showDetail()
				}
			}
		}

		if event.Key() == tcell.KeyEsc {
			if inputField.HasFocus() {
				rootPage.SwitchToPage("sandPage")
				app.SetFocus(inputField)
				for _, t := range timers {
					go t.PH(t.G, t.RF)
				}
				currentView = "sandPage"
			}
		}

		if event.Key() == tcell.KeyEsc {
			if createTimerView.HasFocus() || rootPage.GetPage("sandPage").HasFocus() {
				showList()
			}
		}

		if event.Rune() == 'c' {
			if timerList.HasFocus() {
				if currentView == "listView" {
					showDetail()
					return nil
				}
			}
		}

		if event.Rune() == 'p' {
			if timerList.HasFocus() {
				if currentView == "timerOptionsView" {
					selectedTimer := timers[timerList.GetCurrentItem()]
					pauseTimer(selectedTimer.N, showToast)
				}
			}
		}

		if event.Rune() == 'r' {
			if timerList.HasFocus() {
				if currentView == "timerOptionsView" {
					selectedTimer := timers[timerList.GetCurrentItem()]
					resumeTimer(selectedTimer.N, showToast, showTimerComplete)
				}
			}
		}

		if event.Rune() == 'q' {
			app.Stop()
			fmt.Println("createtimer", createTimerView.HasFocus())
			fmt.Println("inputfield", inputField.HasFocus())
			fmt.Println("list", timerList.HasFocus())
			fmt.Println("grid", grid.HasFocus())
			fmt.Println("focus", app.GetFocus())
			fmt.Println(len(timers))
			fmt.Println(currentView)

			fmt.Println(rootPage.GetFrontPage())
		}

		if event.Rune() == 's' {
			if timerList.HasFocus() {
				selectedTimer := timers[timerList.GetCurrentItem()]
				if currentView == "listView" {
					stopTimer(selectedTimer.N, showToast)
				} else if currentView == "timerOptionsView" {
					showTimeRemaining(selectedTimer.N, showToast)
				}
			}
		}

		return event
	})

	if err := app.SetRoot(rootPage, true).SetFocus(inputField).Run(); err != nil {
		panic("hi")
		// panic(err)
	}
}
