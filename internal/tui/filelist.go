package tui

import (
	"fmt"
	"github.com/daniloqueiroz/dfm/internal/view"
	"github.com/gdamore/tcell"
	"github.com/inhies/go-bytesize"
	"github.com/rivo/tview"
	"sort"
)

const (
	DIR_FMT  = "[blue::b]"
	FILE_FMT = "[white::]"
)

type filelist struct {
	elem *tview.List
}

func (f filelist) update(items []view.FileItem) {
	cleanNames := make(map[string]string, len(items))
	f.elem.Clear()
	var dirs []string
	var files []string
	for _, item := range items {
		if item.IsDir {
			formatted := fmt.Sprintf("%s%s", DIR_FMT, item.Name)
			dirs = append(dirs, formatted)
			cleanNames[formatted] = item.Name
		} else {
			_, _, width, _ := f.elem.GetRect()
			padding := width - (len(item.Name) + 2)
			size := bytesize.New(float64(item.Size))
			formatted := fmt.Sprintf("%s%s%*s", FILE_FMT, item.Name, padding, size)
			files = append(files, formatted)
			cleanNames[formatted] = item.Name
		}
	}
	sort.Strings(dirs)
	for _, dir := range dirs {
		f.elem.AddItem(dir, cleanNames[dir], 0, func() {})
	}
	sort.Strings(files)
	for _, file := range files {
		f.elem.AddItem(file, cleanNames[file], 0, func() {})
	}
}

func (f filelist) registerKeyHandlers(evChan chan interface{}) {
	// TODO fuzzy search if user start typing ?
	// 	"github.com/sahilm/fuzzy"
	//f.elem.FindItems()
	f.elem.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		selectedPos := f.elem.GetCurrentItem()
		_, selectedName := f.elem.GetItemText(selectedPos)
		switch event.Rune() {
		case '+':
			evChan <- view.FileListItemSelected{
				Pos:           selectedPos,
				Name:          selectedName,
				SelectionMode: view.AddSelectionList,
			}
		case '-':
			evChan <- view.FileListItemSelected{
				Pos:           selectedPos,
				Name:          selectedName,
				SelectionMode: view.RemoveSelectionList,
			}
		}
		return event
	})
	f.elem.SetChangedFunc(func(pos int, _ string, name string, _ rune) {
		evChan <- view.FileListItemHover{
			Pos:  pos,
			Name: name,
		}
	})
	f.elem.SetSelectedFunc(func(pos int, _ string, name string, _ rune) {
		evChan <- view.FileListItemSelected{
			Pos:           pos,
			Name:          name,
			SelectionMode: view.Open,
		}
	})
}

func newFilelist() *filelist {
	list := tview.NewList()
	list.SetTitleAlign(tview.AlignCenter)
	list.SetHighlightFullLine(true)
	list.ShowSecondaryText(false)
	return &filelist{
		elem: list,
	}
}
