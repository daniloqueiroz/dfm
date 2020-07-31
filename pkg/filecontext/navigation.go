package filecontext

import (
	"github.com/daniloqueiroz/dfm/pkg/vfs"
	"github.com/gammazero/deque"
)

type Navigation struct {
	prev    deque.Deque
	current vfs.File
	next    deque.Deque
}

func (n *Navigation) Current() vfs.File {
	return n.current
}

func (n *Navigation) SetCurrent(f vfs.File) {
	n.prev.PushFront(n.current)
	n.next.Clear()
	n.current = f
}

func (n *Navigation) Previous() {
	if n.prev.Len() > 0 {
		n.next.PushFront(n.current)
		n.current = n.prev.PopFront().(vfs.File)
	}
}

func (n *Navigation) Next() {
	if n.next.Len() > 0 {
		n.prev.PushFront(n.current)
		n.current = n.next.PopFront().(vfs.File)
	}
}

func NewNavigation(current vfs.File) *Navigation {
	return &Navigation{
		current: current,
	}
}
