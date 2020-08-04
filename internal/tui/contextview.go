package tui

import (
	"fmt"
	"github.com/daniloqueiroz/dfm/internal/view"
	"github.com/inhies/go-bytesize"
	"github.com/rivo/tview"
	"strings"
	"time"
)

type contextpane struct {
	elem *tview.TextView
}

func (c contextpane) update(details interface{}) {
	text := ""
	switch d := details.(type) {
	case view.FileDetails:
		text = fmt.Sprintf(
			"%s\n%s\n%s",
			d.ModTime.Format(time.UnixDate),
			d.Mode.String(),
			bytesize.New(float64(d.Size)),
		)
		//c.elem.SetTitle("File Details")
	case []view.FileItem:
		var b strings.Builder
		for _, file := range d {
			b.WriteString(fmt.Sprintf("%s\n", file.FullPath))
		}
		text = b.String()
		//c.elem.SetTitle("Selected Files")
	}
	if text != "" {
		c.elem.SetText(text)
	} else {
		c.elem.Clear()
	}
}

func newContextPane() *contextpane {
	d := tview.NewTextView()
	d.SetBorder(true)
	//d.SetTitle("File Details")
	return &contextpane{
		elem: d,
	}
}
