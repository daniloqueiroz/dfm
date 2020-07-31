package cui

import (
	"fmt"
	"github.com/daniloqueiroz/dfm/internal/view"
	"github.com/rivo/tview"
)

type statusbar struct {
	elem *tview.TextView
}

func (s statusbar) update(info view.Status) {
	counters := fmt.Sprintf("# files: %d", info.FilesCount)
	if info.SelectedFilesCount > 0 {
		counters = fmt.Sprintf("# files selected: %d | %s", info.SelectedFilesCount, counters)
	}
	_, _, width, _ := s.elem.GetRect()
	padding := width - (len(counters) + 4 + len(info.Message))
	msg := fmt.Sprintf("[%d] %s%*s", info.Context, info.Message, padding, counters)
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
