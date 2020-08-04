package tui

import (
	"fmt"
	"github.com/daniloqueiroz/dfm/internal/view"
	"github.com/rivo/tview"
)

type statusbar struct {
	elem *tview.TextView
}

func (s statusbar) update(info view.Status) {
	counters := fmt.Sprintf("# files selected: %d | # files: %d", info.SelectedFilesCount, info.FilesCount)
	_, _, width, _ := s.elem.GetRect()
	prefix := fmt.Sprintf("[ %d ]", info.Context)
	padding := width - (len(prefix) + len(counters))
	msg := fmt.Sprintf("%s%*s", prefix, padding, counters)
	s.elem.SetText(msg)
}

func newStatusBar() *statusbar {
	s := tview.NewTextView()
	s.SetBorder(false)
	//s.SetTextColor(tcell.ColorDarkGray)
	//s.SetBackgroundColor(tcell.ColorWhite)
	s.SetText("statusbar")
	return &statusbar{
		elem: s,
	}
}
