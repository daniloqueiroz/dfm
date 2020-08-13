package view

type ScreenSizeChanged struct{}
type ToggleCommandMode struct {
	Prefix rune
}
type Command struct {
	Cmdline string
}
type ToggleHiddenFilesVisibility struct{}
type ToggleFileSelectionView struct{}
type NavNext struct{}
type NavPrev struct{}
type FileListItemHover struct {
	Pos  int
	Name string
}
type FileSelectionType string

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
