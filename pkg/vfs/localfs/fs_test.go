package localfs

import (
	"fmt"
	"github.com/daniloqueiroz/dfm/pkg/vfs"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func randomName() string {
	return fmt.Sprintf("gotest-%d", time.Now().Unix())
}

func tmpFS() vfs.FileSystem {
	fs, err := NewFS("/tmp")
	if err != nil {
		panic("unable to create /tmp fs")
	}
	return fs
}

func isFileThere(path string, files []vfs.File) bool {
	found := false
	for _, file := range files {
		if file.Path() == path {
			found = true
		}
	}
	return found
}

func TestFS_createDirListGetDeleteDir(t *testing.T) {
	fs := tmpFS()
	dirname := randomName()

	file, err := fs.CreateDir(dirname)
	assert.Nil(t, err)

	files, err := fs.List()
	assert.Nil(t, err)
	assert.True(t, isFileThere(file.Path(), files))

	getFile, err := fs.GetFile(dirname)
	assert.Nil(t, err)
	assert.Equal(t, file.Path(), getFile.Path())

	err = fs.Delete(file)
	assert.Nil(t, err)
	assert.True(t, isFileThere(file.Path(), files))
}
