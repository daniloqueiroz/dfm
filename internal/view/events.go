package view

type ViewSizeChanged struct{}

type NavNext struct{}
type NavPrev struct{}
type FileListItemHover struct {
	Pos  int
	Name string
}
type FileListItemSelected struct {
	Pos  int
	Name string
}
