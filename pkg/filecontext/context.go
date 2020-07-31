package filecontext

import (
	"github.com/daniloqueiroz/dfm/pkg/vfs"
	"github.com/google/logger"
)

type FileContext struct {
	nav *Navigation
	fs  vfs.FileSystem
}

func (c FileContext) GetCWD() vfs.File {
	return c.fs.GetCWD()
}

func (c FileContext) ListFiles() ([]vfs.File, error) {
	return c.fs.List()
}

func (c FileContext) Next() {
	c.nav.Next()
	err := c.fs.SetCWD(c.nav.current)
	if err != nil {
		logger.Errorf("Error setting fs cwd: %v", err)
		c.nav.Previous()
	}
}

func (c FileContext) Previous() {
	c.nav.Previous()
	err := c.fs.SetCWD(c.nav.current)
	if err != nil {
		logger.Errorf("Error setting fs cwd: %v", err)
		c.nav.Next()
	}
}

func (c FileContext) CD(name string) error {
	file, err := c.fs.GetFile(name)
	if err != nil {
		return err
	}

	err = c.fs.SetCWD(file)
	if err != nil {
		return nil
	}

	c.nav.SetCurrent(file)
	return nil
}

func NewFileContext(fs vfs.FileSystem) *FileContext {
	root := fs.GetCWD()
	return &FileContext{
		nav: NewNavigation(root),
		fs:  fs,
	}
}
