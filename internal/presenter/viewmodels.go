package presenter

import (
	"github.com/daniloqueiroz/dfm/internal/view"
)

type ViewElement uint8

const (
	Navigation AppMode = "navigation"
	Command    AppMode = "command"
)

type viewConfig struct {
	mode          AppMode
	hideHidden    bool
	showSelection bool // context mode
}

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
	location          string
	fileList          []view.FileItem
	status            view.Status
	fileDetail        *view.FileDetails
	selectedList      []view.FileItem
	commandBarContent string
	toRefresh         ViewElement
}
