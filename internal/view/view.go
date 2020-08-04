package view

import (
	"os"
	"time"
)

// View data models
type FileItem struct {
	Name     string
	FullPath string
	IsDir    bool
	Size     int64
}

type Status struct {
	Context            int
	FilesCount         int
	SelectedFilesCount int
}

type FileDetails struct {
	Size    int64
	Mode    os.FileMode
	ModTime time.Time
}

type View interface {
	Show(eventHandler func(interface{}))
	SetLocationBar(path string)
	SetFileList(items []FileItem)
	SetStatusMessage(info Status)
	SetContextDetails(details interface{})
	Quit()
}
