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
	ALL = LocationBar | FileListView | ContextView | StatusBar
)

type viewData struct {
	location     string
	fileList     []view.FileItem
	status       view.Status
	fileDetail   *view.FileDetails
	selectedList []view.FileItem
	toRefresh    ViewElement
}

func refresh(v view.View, data viewData, cfg config) {
	if data.toRefresh&LocationBar != 0 {
		v.SetLocationBar(data.location)
	}
	if data.toRefresh&FileListView != 0 {
		v.SetFileList(data.fileList)
	}
	if data.toRefresh&ContextView != 0 {
		if cfg.showSelection {
			v.SetContextDetails(data.selectedList)
		} else {
			if data.fileDetail != nil {
				v.SetContextDetails(*data.fileDetail)
			} else {
				v.SetContextDetails(nil)
			}
		}
	}
	if data.toRefresh&StatusBar != 0 {
		v.SetStatusMessage(data.status)
	}
}
