package tui

import (
	"fmt"
	"github.com/daniloqueiroz/dfm/internal/view"
	"github.com/gdamore/tcell"
	"github.com/rivo/tview"
)

type statusbar struct {
	elem *tview.TextView
}

func (s statusbar) update(info view.Status) {
	counters := fmt.Sprintf(" # files: %d | # files selected: %d ", info.FilesCount, info.SelectedFilesCount)
	_, _, width, _ := s.elem.GetRect()
	prefix := fmt.Sprintf("[ %d ]", info.Context)
	padding := width - len(prefix)
	msg := fmt.Sprintf("%s%*s", prefix, padding, counters)
	s.elem.SetText(msg)
}

func (s statusbar) fade() {
	s.elem.SetTextColor(tcell.ColorWhite)
	s.elem.SetBackgroundColor(tcell.ColorBlack)
}

func (s statusbar) highlight() {
	s.elem.SetTextColor(tcell.ColorBlack)
	s.elem.SetBackgroundColor(tcell.ColorWhite)
}

func newStatusBar() *statusbar {
	s := tview.NewTextView()
	s.SetBorder(false)
	s.SetTextColor(tcell.ColorBlack)
	s.SetBackgroundColor(tcell.ColorWhite)
	s.SetText("statusbar")
	return &statusbar{
		elem: s,
	}
}
