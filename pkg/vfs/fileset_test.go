package vfs

import (
	"github.com/stretchr/testify/assert"
	"os"
	"path/filepath"
	"testing"
)

type dummyFile struct {
	path string
}

func (d *dummyFile) Parent() File {
	return nil
}
func (d *dummyFile) Exists() bool {
	return true
}

func (d *dummyFile) IsHidden() bool {
	return false
}

func (d *dummyFile) Name() string {
	return filepath.Base(d.path)
}

func (d *dummyFile) Path() string {
	return d.path
}

func (d *dummyFile) Stats() (os.FileInfo, error) {
	return nil, nil
}

func TestNewFileSet_AddContains(t *testing.T) {
	file1 := &dummyFile{
		path: "/tmp/file1",
	}

	set := NewFileSet()
	assert.False(t, set.Contains(file1))

	set.Add(file1)
	assert.True(t, set.Contains(file1))
}

func TestNewFileSet_AddNoDuplicates(t *testing.T) {
	file1 := &dummyFile{
		path: "/tmp/file1",
	}
	file2 := &dummyFile{
		path: "/tmp/file2",
	}

	set := NewFileSet()
	assert.Equal(t, 0, set.Len())

	set.Add(file1)
	assert.Equal(t, 1, set.Len())

	set.Add(file2)
	assert.Equal(t, 2, set.Len())

	copyFile1 := &dummyFile{
		path: "/tmp/file1",
	}
	set.Add(copyFile1)
	assert.Equal(t, 2, set.Len())
}

func TestNewFileSet_Clear(t *testing.T) {
	file1 := &dummyFile{
		path: "/tmp/file1",
	}
	file2 := &dummyFile{
		path: "/tmp/file2",
	}

	set := NewFileSet()
	set.Add(file1)
	set.Add(file2)
	assert.Equal(t, 2, set.Len())

	set.Clear()
	assert.Equal(t, 0, set.Len())
}

func TestNewFileSet_Iterate(t *testing.T) {
	file1 := &dummyFile{
		path: "/tmp/file1",
	}
	file2 := &dummyFile{
		path: "/tmp/file2",
	}

	set := NewFileSet()
	set.Add(file1)
	set.Add(file2)
	assert.Equal(t, 2, set.Len())

	itr := set.Iterator()
	assert.NotNil(t, itr)
	assert.Contains(t, itr, file1)
	assert.Contains(t, itr, file2)
}

func TestNewFileSet_AddAndRemove(t *testing.T) {
	file1 := &dummyFile{
		path: "/tmp/file1",
	}

	set := NewFileSet()
	assert.Equal(t, 0, set.Len())

	set.Add(file1)
	assert.Equal(t, 1, set.Len())

	set.Remove(file1)
	assert.Equal(t, 0, set.Len())
	assert.False(t, set.Contains(file1))
}
