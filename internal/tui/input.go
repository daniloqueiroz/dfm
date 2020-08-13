package tui

import (
	"github.com/daniloqueiroz/dfm/internal/view"
	"github.com/gdamore/tcell"
	"github.com/google/logger"
	"github.com/rivo/tview"
)

type input struct {
	elem *tview.InputField
}

func (i *input) Clear() {
	logger.Infof("clear")
	i.elem.SetText("")
}

func (i *input) setText(prefix string) {
	i.elem.SetText(prefix)
}

func (i *input) registerKeyHandlers(evChan chan interface{}) {
	i.elem.SetDoneFunc(func(key tcell.Key) {
		if key == tcell.KeyEnter {
			evChan <- view.Command{Cmdline: i.elem.GetText()}
		}
	})
}

func newCommandBar() *input {
	i := tview.NewInputField()
	i.SetBorder(false)
	i.SetFieldBackgroundColor(tcell.ColorBlack)
	i.SetFieldTextColor(tcell.ColorWhite)
	return &input{
		elem: i,
	}
}
