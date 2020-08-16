package tui

import (
	"fmt"
	"github.com/daniloqueiroz/dfm/internal/view"
	"github.com/gdamore/tcell"
	"github.com/rivo/tview"
	"strings"
)

type statusbar struct {
	elem *tview.TextView
}

func (s statusbar) update(info view.Status) {
	counters := fmt.Sprintf(" # files: %d | # files selected: %d ", info.FilesCount, info.SelectedFilesCount)
	_, _, width, _ := s.elem.GetRect()
	var prefixBuilder strings.Builder
	prefixBuilder.WriteString(" [")
	for idx := 0; idx < info.ContextCount; idx++ {
		if idx == info.ActiveContext {
			prefixBuilder.WriteString(fmt.Sprintf(" [::bu]%d[::-]", idx))
		} else {
			prefixBuilder.WriteString(fmt.Sprintf(" %d", idx))
		}
	}
	prefixBuilder.WriteString(" ]")
	prefix := prefixBuilder.String()
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
	s.SetDynamicColors(true)
	return &statusbar{
		elem: s,
	}
}
