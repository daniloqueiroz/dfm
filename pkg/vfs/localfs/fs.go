package localfs

import (
	"fmt"
	"github.com/bluele/gcache"
	"github.com/daniloqueiroz/dfm/pkg/vfs"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"
)

type localFS struct {
	cache gcache.Cache // TODO implement caching strategy
	cwd   vfs.File
}

func (fs *localFS) stats() os.FileInfo {
	stats, err := fs.cwd.Stats()
	if err != nil {
		panic(err)
	}
	return stats
}

func (fs *localFS) GetCWD() vfs.File {
	return fs.cwd
}

func (fs *localFS) SetCWD(file vfs.File) error {
	stats, err := file.Stats()
	if err != nil {
		return fmt.Errorf("fs: %s can't be cwd: %w", file.Path(), err)
	} else if !stats.IsDir() {
		return fmt.Errorf("fs: %s can't be cwd: it isn't a dir", file.Path())
	}

	fs.cwd = file
	return nil
}

func (fs *localFS) List() ([]vfs.File, error) {
	var files []vfs.File
	infos, err := ioutil.ReadDir(fs.cwd.Path())
	if err != nil {
		return nil, err
	}
	for _, info := range infos {
		path := filepath.Join(fs.cwd.Path(), info.Name())
		files = append(files, newFile(path, true, info))
	}
	return files, nil
}

func (fs *localFS) GetFile(name string) (vfs.File, error) {
	files, err := fs.List()
	if err != nil {
		return nil, err
	}
	for _, f := range files {
		if f.Name() == name {
			return f, nil
		}
	}
	return nil, fmt.Errorf("fs: file %s not found on dir %s", name, fs.GetCWD().Name())
}

func (fs *localFS) CreateFile(name string) (vfs.File, error) {
	panic("implement me")
}

func (fs *localFS) CreateDir(name string) (vfs.File, error) {
	if filepath.IsAbs(name) {
		// TODO error
	}
	path := filepath.Join(fs.cwd.Path(), name)
	err := os.Mkdir(path, fs.stats().Mode().Perm())
	if err != nil {
		return nil, err
	}
	return NewFile(path)
}

func (fs *localFS) CreateDirs(name string) error {
	panic("implement me")
}

func (fs *localFS) Copy(file vfs.File) error {
	panic("implement me")
}

func (fs *localFS) Move(file vfs.File) error {
	panic("implement me")
}

func (fs *localFS) Delete(file vfs.File) error {
	dir, _ := filepath.Split(file.Path())
	if fs.cwd.Path() != filepath.Clean(dir) {
		return fmt.Errorf("fs: can't remove file: file %s doesn't belong to %s", file.Path(), fs.cwd.Path())
		// TODO handle corner cases
		// file not in dir, non empty subdirs
	}
	return os.Remove(file.Path())
}

func (fs *localFS) Rename(file vfs.File) error {
	panic("implement me")
}

func NewFS(startDir string) (vfs.FileSystem, error) {
	file, err := NewFile(startDir)
	if err != nil {
		return nil, fmt.Errorf("fs: unable to initialize fs: %w", err)
	}

	stats, err := file.Stats()
	if err != nil {
		return nil, fmt.Errorf("fs: %s can be start dir: %w", startDir, err)
	} else if !stats.IsDir() {
		return nil, fmt.Errorf("fs: %s can be start dir: it isn't a dir", startDir)
	} else {
		return &localFS{
			// TODO move cache parameters to config
			cache: gcache.New(20).LRU().Expiration(60 * time.Second).Build(),
			cwd:   file,
		}, nil
	}
}
