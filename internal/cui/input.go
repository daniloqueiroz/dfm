package cui

import (
	"github.com/gdamore/tcell"
	"github.com/rivo/tview"
)

type input struct {
	elem *tview.InputField
}

func newCommandBar() *input {
	i := tview.NewInputField()
	i.SetBorder(true)
	i.SetPlaceholder("placehold")
	i.SetBackgroundColor(tcell.ColorBlack)
	i.SetFieldBackgroundColor(tcell.ColorBlack)
	i.SetFieldTextColor(tcell.ColorWhite)
	i.SetPlaceholderTextColor(tcell.ColorGray)
	return &input{
		elem: i,
	}
}
