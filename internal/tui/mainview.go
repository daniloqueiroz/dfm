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
}

func (w *Window) keyHandler(event *tcell.EventKey) *tcell.EventKey {
	result := event
	switch event.Key() {
	case tcell.KeyRight:
		w.evChan <- view.NavNext{}
		return nil
	case tcell.KeyLeft:
		w.evChan <- view.NavPrev{}
		return nil
	}

	switch event.Rune() {
	case 'q':
		w.evChan <- view.Quit{}
	case 'h':
		w.evChan <- view.ToggleHiddenFilesVisibility{}
		return nil
	case 'S':
		w.evChan <- view.ToggleFileSelectionView{}
		return nil
	}

	return result
}

func (w *Window) registerKeyHandlers() {
	w.app.SetInputCapture(w.keyHandler)
	w.flView.registerKeyHandlers(w.evChan)
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

func (w *Window) Show(handler func(interface{})) {
	defer internal.OnPanic("Window:Show")
	go func() {
		for ev := range w.evChan {
			handler(ev)
		}
	}()

	grid := tview.NewGrid()
	grid.AddItem(w.lcBar.elem, 0, 0, 1, 4, 0, 0, false)
	grid.AddItem(w.lcPane.elem, 1, 0, 7, 3, 0, 0, false)
	grid.AddItem(w.flView.elem, 1, 1, 7, 2, 0, 0, true)
	grid.AddItem(w.ctxView.elem, 1, 3, 7, 1, 0, 0, false)
	grid.AddItem(w.stBar.elem, 8, 0, 1, 4, 0, 0, false)
	//grid.AddItem(cmdBar, 8,0,1,4,0,0,false)
	w.app.SetRoot(grid, true).SetFocus(w.flView.elem)

	w.app.SetAfterDrawFunc(w.afterDraw)
	w.registerKeyHandlers()

	logger.Infof("Starting app loop")
	if err := w.app.Run(); err != nil {
		panic(err)
	}
}

func (w *Window) Quit() {
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
	}
}
