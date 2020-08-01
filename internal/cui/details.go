package cui

import (
	"fmt"
	"github.com/daniloqueiroz/dfm/internal/view"
	"github.com/google/logger"
	"github.com/inhies/go-bytesize"
	"github.com/rivo/tview"
	"time"
)

type contextpane struct {
	elem *tview.TextView
}

func (c contextpane) update(details interface{}) {
	switch d := details.(type) {
	case view.FileDetails:
		logger.Infof("Details received: %+v", d)
		text := fmt.Sprintf(
			"%s\n%s\n%s",
			d.ModTime.Format(time.UnixDate),
			d.Mode.String(),
			bytesize.New(float64(d.Size)),
		)
		c.elem.SetText(text)
	default:
		c.elem.Clear()
	}
}

func newContextPane() *contextpane {
	d := tview.NewTextView()
	d.SetBorder(true)
	return &contextpane{
		elem: d,
	}
}
