package tui

import (
	"fmt"
	"github.com/daniloqueiroz/dfm/internal"
	"github.com/daniloqueiroz/dfm/internal/view"
	"github.com/gdamore/tcell"
	"github.com/google/logger"
	"github.com/rivo/tview"
)

type Window struct {
	evChan     chan interface{}
	app        *tview.Application
	lcBar      *locationbar
	lcPane     *locationpane
	flView     *filelist
	ctxView    *contextpane
	stBar      *statusbar
	cmdBar     *input
	screenSize string
	isFocusFl  bool
}

func (w *Window) keyHandler(event *tcell.EventKey) *tcell.EventKey {
	result := event
	switch event.Key() {
	case tcell.KeyRight:
		if w.isFocusFl {
			w.evChan <- view.NavNext{}
			return nil
		}
	case tcell.KeyLeft:
		if w.isFocusFl {
			w.evChan <- view.NavPrev{}
			return nil
		}
	case tcell.KeyESC:
		w.evChan <- view.ToggleCommandMode{
			Prefix: 0,
		}
		return nil
	}

	if w.isFocusFl {
		switch event.Rune() {
		case ':':
			w.evChan <- view.ToggleCommandMode{
				Prefix: ':',
			}
		case '/':
			w.evChan <- view.ToggleCommandMode{
				Prefix: '/',
			}
		case 'S':
			w.evChan <- view.ToggleFileSelectionView{}
			return nil
		}
	}

	return result
}

func (w *Window) registerKeyHandlers() {
	w.app.SetInputCapture(w.keyHandler)
	w.flView.registerKeyHandlers(w.evChan)
	w.cmdBar.registerKeyHandlers(w.evChan)
}

func (w *Window) afterDraw(s tcell.Screen) {
	x, y := s.Size()
	screenSize := fmt.Sprintf("%d;%d", x, y)
	if w.screenSize == "" || w.screenSize != screenSize {
		logger.Infof("Screen size changed: %s", screenSize)
		w.screenSize = screenSize
		w.evChan <- view.ScreenSizeChanged{}
	}
}

func (w *Window) SetLocationBar(path string) {
	w.app.QueueUpdateDraw(func() {
		w.lcBar.elem.SetText(path)
	})
}

func (w *Window) SetFileList(items []view.FileItem) {
	w.app.QueueUpdateDraw(func() {
		w.flView.update(items)
	})
}

func (w *Window) SetStatusMessage(info view.Status) {
	w.app.QueueUpdateDraw(func() {
		w.stBar.update(info)
	})
}

func (w *Window) SetContextDetails(details interface{}) {
	w.app.QueueUpdateDraw(func() {
		w.ctxView.update(details)
	})
}

func (w *Window) OnEvent(handler func(interface{})) {
	go func() {
		for ev := range w.evChan {
			handler(ev)
		}
	}()
}

func (w *Window) Show() {
	defer internal.OnPanic("Window:Show")
	emptyBox := tview.NewBox()
	emptyBox.SetBorder(false)
	contentLayout := tview.NewFlex()
	contentLayout.AddItem(w.flView.elem, 0, 3, false)
	contentLayout.AddItem(emptyBox, 1, 0, false)
	contentLayout.AddItem(w.ctxView.elem, 0, 1, false)
	contentLayout.AddItem(w.lcPane.elem, 0, 0, false)

	windowLayout := tview.NewFlex()
	windowLayout.SetDirection(tview.FlexRow)
	windowLayout.AddItem(w.lcBar.elem, 1, 0, false)
	windowLayout.AddItem(emptyBox, 1, 0, false)
	windowLayout.AddItem(contentLayout, 0, 1, false)
	windowLayout.AddItem(w.stBar.elem, 1, 0, false)
	windowLayout.AddItem(w.cmdBar.elem, 1, 0, false)
	w.app.SetRoot(windowLayout, true).SetFocus(w.flView.elem)

	w.app.SetAfterDrawFunc(w.afterDraw)
	w.registerKeyHandlers()

	logger.Infof("Starting app loop")
	if err := w.app.Run(); err != nil {
		panic(err)
	}
}

func (w *Window) SetCommandBar(prefix string) {
	w.cmdBar.setText(prefix)
}

func (w *Window) FocusCommandBar() {
	w.isFocusFl = false
	w.app.SetFocus(w.cmdBar.elem)
	w.stBar.fade()
	w.app.Draw()
}

func (w *Window) FocusFileList() {
	w.isFocusFl = true
	w.cmdBar.Clear()
	w.app.SetFocus(w.flView.elem)
	w.stBar.highlight()
	w.app.Draw()
}

func (w *Window) Quit() {
	close(w.evChan)
	w.app.Stop()
}

func NewWindow() view.View {
	return &Window{
		app:        tview.NewApplication(),
		lcBar:      newLocationBar(),
		lcPane:     newLocationPane(),
		flView:     newFilelist(),
		ctxView:    newContextPane(),
		stBar:      newStatusBar(),
		cmdBar:     newCommandBar(),
		screenSize: "",
		evChan:     make(chan interface{}),
		isFocusFl:  true,
	}
}
