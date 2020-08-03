package localfs

import (
	"fmt"
	"github.com/daniloqueiroz/dfm/pkg/vfs"
	"os"
	"path/filepath"
	"strings"
)

type localFile struct {
	name     string
	fullpath string
	exists   bool
	isHidden bool
	info     os.FileInfo
}

func (f *localFile) Parent() vfs.File {
	if f.fullpath == "/" {
		return nil
	}

	var parentPath string
	stats, _ := f.Stats()
	if stats != nil && stats.IsDir() {
		parentPath = filepath.Clean(filepath.Join(f.fullpath, ".."))
	} else {
		parentPath = filepath.Dir(f.fullpath)
	}

	file, _ := NewFile(parentPath)
	return file
}

func (f *localFile) Exists() bool {
	return f.exists
}

func (f *localFile) IsHidden() bool {
	return f.isHidden
}

func (f *localFile) Name() string {
	return f.name
}

func (f *localFile) Path() string {
	return f.fullpath
}

func (f *localFile) Stats() (os.FileInfo, error) {
	if f.exists {
		return f.info, nil
	} else {
		return nil, fmt.Errorf("file: file %s doesn't exists", f.fullpath)
	}
}

func newFile(name string, exists bool, stats os.FileInfo) vfs.File {
	basename := filepath.Base(name)

	return &localFile{
		name:     basename,
		fullpath: name,
		exists:   exists,
		info:     stats,
		isHidden: strings.HasPrefix(basename, "."),
	}
}

func NewFile(name string) (vfs.File, error) {
	if !filepath.IsAbs(name) {
		return nil, fmt.Errorf("file: unable to create file, %s isn't a absolute path", name)
	}
	stats, err := os.Stat(name)
	return newFile(name, err == nil, stats), nil
}
