package tui

import (
	"github.com/daniloqueiroz/dfm/internal/view"
	"github.com/gdamore/tcell"
	"github.com/rivo/tview"
	"strconv"
	"strings"
)

const (
	QuitCommand                 = ":quit"
	QuitCommandShort            = ":q"
	ToggleCommand               = ":toggle"
	HiddenToggleCommandShort    = ":h"
	SelectionToggleCommandShort = ":s"
	Context                     = ":context"
	NewContextShort             = ":n"
	CloseContextShort           = ":w"
)

type commandBar struct {
	elem *tview.InputField
}

func (cb *commandBar) Clear() {
	cb.elem.SetText("")
}

func (cb *commandBar) setText(prefix string) {
	cb.elem.SetText(prefix)
}

func (cb *commandBar) registerKeyHandlers(evChan chan interface{}) {
	cb.elem.SetDoneFunc(func(key tcell.Key) {
		if key == tcell.KeyEnter {
			cmdline := cb.elem.GetText()
			if cb.handleCommand(cmdline, evChan) {
				return
			}
		}
	})
}

func (cb *commandBar) handleCommand(cmdline string, evChan chan interface{}) bool {
	tokens := strings.Split(cmdline, " ")
	cmd := tokens[0]
	params := strings.TrimSpace(cmdline[len(cmd):])
	switch cmd {
	case QuitCommand, QuitCommandShort:
		evChan <- view.QuitEvent{}
	case NewContextShort:
		evChan <- view.CreateContext{}
	case CloseContextShort:
		evChan <- view.CloseContext{Index: -1}
	case SelectionToggleCommandShort:
		evChan <- view.ToggleSelectionView{}
	case HiddenToggleCommandShort:
		evChan <- view.ToggleHiddenFiles{}
	case ToggleCommand:
		if params == "selection" || params == "details" {
			evChan <- view.ToggleSelectionView{}
		} else if params == "hidden" {
			evChan <- view.ToggleHiddenFiles{}
		}
	case Context:
		if len(tokens) < 2 {
			return true
		} else {
			subcommand := tokens[1]
			switch subcommand {
			case "new":
				if len(tokens) > 2 {
					params = tokens[2]
				} else {
					params = ""
				}
				evChan <- view.CreateContext{
					BaseDir: params,
				}
			case "close":
				idx := -1
				if len(tokens) > 2 {
					var err error
					idx, err = strconv.Atoi(tokens[2])
					if err != nil {
						return true
					}
				}
				evChan <- view.CloseContext{
					Index: idx,
				}
			}
		}
	}
	return false
}

func newCommandBar() *commandBar {
	i := tview.NewInputField()
	i.SetBorder(false)
	i.SetFieldBackgroundColor(tcell.ColorBlack)
	i.SetFieldTextColor(tcell.ColorWhite)
	return &commandBar{
		elem: i,
	}
}
