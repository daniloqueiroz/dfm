package view

type ScreenSizeChanged struct{}
type ToggleCommandMode struct {
	Prefix rune
}
type Command struct {
	Cmdline string
}
type NavNext struct{}
type NavPrev struct{}
type FileListItemHover struct {
	Pos  int
	Name string
}
type FileSelectionType string
type SwitchContext struct {
	Index int
}

const (
	Open                FileSelectionType = "open"
	AddSelectionList    FileSelectionType = "add_selection"
	RemoveSelectionList FileSelectionType = "remove_selection"
)

type FileListItemSelected struct {
	Pos           int
	Name          string
	SelectionMode FileSelectionType
}
