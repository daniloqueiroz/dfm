package view

type ViewSizeChanged struct{}

type NavNext struct{}
type NavPrev struct{}
type FileListItemSelected struct {
	Pos  int
	Name string
}
