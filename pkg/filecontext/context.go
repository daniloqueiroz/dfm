package filecontext

import (
	"github.com/daniloqueiroz/dfm/pkg/vfs"
	"github.com/google/logger"
	"os"
)

type FileContext struct {
	nav *Navigation
	fs  vfs.FileSystem
}

func (fc FileContext) GetCWD() vfs.File {
	return fc.fs.GetCWD()
}

func (fc FileContext) ListFiles() ([]vfs.File, error) {
	return fc.fs.List()
}

func (fc FileContext) Next() {
	fc.nav.Next()
	err := fc.fs.SetCWD(fc.nav.current)
	if err != nil {
		logger.Errorf("Error setting fs cwd: %v", err)
		fc.nav.Previous()
	}
}

func (fc FileContext) Previous() {
	fc.nav.Previous()
	err := fc.fs.SetCWD(fc.nav.current)
	if err != nil {
		logger.Errorf("Error setting fs cwd: %v", err)
		fc.nav.Next()
	}
}

func (fc FileContext) CD(name string) error {
	var file vfs.File
	var err error

	if name == ".." { // TODO is this the right place for this?
		file = fc.GetCWD().Parent()
	} else {
		file, err = fc.fs.GetFile(name)
		if err != nil {
			return err
		}
	}

	err = fc.fs.SetCWD(file)
	if err != nil {
		return nil
	}

	fc.nav.SetCurrent(file)
	return nil
}

func (fc FileContext) Stats(name string) (os.FileInfo, error) {
	if name == ".." {
		return nil, nil
	} else {
		file, err := fc.fs.GetFile(name)
		if err != nil {
			return nil, err
		} else {
			return file.Stats()
		}
	}
}

func (fc FileContext) GetFile(name string) (vfs.File, error) {
	return fc.fs.GetFile(name)
}

func NewFileContext(fs vfs.FileSystem) *FileContext {
	root := fs.GetCWD()
	return &FileContext{
		nav: NewNavigation(root),
		fs:  fs,
	}
}
