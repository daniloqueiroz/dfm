package presenter

import (
	"github.com/daniloqueiroz/dfm/internal/view"
)

type ViewElement uint8

const (
	LocationBar ViewElement = 1 << iota
	FileListView
	ContextView
	StatusBar
	CommandBar
	Focus
	ALL = LocationBar | FileListView | ContextView | StatusBar | CommandBar | Focus
)

type viewData struct {
	view              view.View
	location          string
	fileList          []view.FileItem
	status            view.Status
	fileDetail        *view.FileDetails
	selectedList      []view.FileItem
	commandBarContent string
	toRefresh         ViewElement
}

func (vd *viewData) start() {
	vd.view.Show()
}
