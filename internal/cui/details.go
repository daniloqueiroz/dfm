package cui

import "github.com/rivo/tview"

type contextpane struct {
	elem *tview.TextView
}

func newContextPane() *contextpane {
	d := tview.NewTextView()
	d.SetBorder(true)
	d.SetText("d")
	return &contextpane{
		elem: d,
	}
}
