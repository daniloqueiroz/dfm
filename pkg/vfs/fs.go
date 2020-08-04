package vfs

import "os"

type FileSet interface {
	Len() int
	Contains(File) bool
	Add(File)
	Remove(File)
	Clear()
	Iterator() []File
}

type File interface {
	Parent() File
	Exists() bool
	IsHidden() bool
	Name() string
	Path() string
	/*
	 * Return the file info if file exists, error if file doesn't exists
	 */
	Stats() (os.FileInfo, error)
}

type FileSystem interface {
	GetCWD() File
	SetCWD(file File) error
	List() ([]File, error)
	GetFile(name string) (File, error)
	CreateFile(name string) (File, error)
	CreateDir(name string) (File, error)
	CreateDirs(name string) error // ?
	Copy(file File) error
	Move(file File) error
	Delete(file File) error
	Rename(file File) error
}
