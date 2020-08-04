package tui

import (
	"github.com/gdamore/tcell"
	"github.com/rivo/tview"
)

type locationbar struct {
	elem *tview.TextView
}

func (l locationbar) update(path string) {
	l.elem.SetText(path)
}

func newLocationBar() *locationbar {
	l := tview.NewTextView()
	l.SetBorder(false)
	l.SetTextColor(tcell.ColorDarkBlue)
	l.SetTextAlign(tview.AlignCenter)
	return &locationbar{
		elem: l,
	}
}

type locationpane struct {
	elem *tview.List
}

func newLocationPane() *locationpane {
	list := tview.NewList()
	list.SetBorder(false)
	list.ShowSecondaryText(false)
	list.AddItem("[::b]root", "", 'r', func() {})
	list.AddItem("[::b]home", "", 'h', func() {})
	return &locationpane{
		elem: list,
	}
}
