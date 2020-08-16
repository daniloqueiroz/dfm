package view

type FileSelectionType string

const (
	Open                FileSelectionType = "open"
	AddSelectionList    FileSelectionType = "add_selection"
	RemoveSelectionList FileSelectionType = "remove_selection"
)

// Events fired by view to presenter
type QuitEvent struct{}
type ScreenSizeChanged struct{}
type ToggleSelectionView struct{}
type ToggleCommandMode struct {
	Prefix rune
}
type ToggleHiddenFiles struct{}
type NavNext struct{}
type NavPrev struct{}
type FileListItemHover struct {
	Pos  int
	Name string
}
type SwitchContext struct {
	Index int
}
type CreateContext struct {
	BaseDir string
}
type CloseContext struct {
	Index int
}
type FileListItemSelected struct {
	Pos           int
	Name          string
	SelectionMode FileSelectionType
}
