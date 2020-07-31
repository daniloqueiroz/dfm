package view

type FileItem struct {
	Name     string
	FullPath string
	IsDir    bool
	Size     int64
}

type Status struct {
	Context            int
	Message            string
	FilesCount         int
	SelectedFilesCount int
}

type View interface {
	Show(eventHandler func(interface{}))
	SetCurrentDir(path string)
	UpdateFileList(items []FileItem)
	SetStatusMessage(info Status)
	ToggleInputBar()
	SetDetails(details string)
}
