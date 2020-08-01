package localfs

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewFile_notAbsolutePath(t *testing.T) {
	_, err := NewFile("test.go")
	assert.Errorf(t, err, "file: unable to create file, test.go isn't a absolute path")
}

func TestNewFile_absolutePathDirExists(t *testing.T) {
	file, err := NewFile("/tmp")
	assert.Nil(t, err)

	assert.True(t, file.Exists())
	assert.Equal(t, "tmp", file.Name())
	assert.Equal(t, "/tmp", file.Path())
	stats, err := file.Stats()
	assert.Nil(t, err)
	assert.NotNil(t, stats)
	assert.True(t, stats.IsDir())
}

func TestNewFile_absolutePathFileExists(t *testing.T) {
	file, err := NewFile("/etc/passwd")
	assert.Nil(t, err)

	assert.True(t, file.Exists())
	assert.Equal(t, "passwd", file.Name())
	assert.Equal(t, "/etc/passwd", file.Path())
	stats, err := file.Stats()
	assert.Nil(t, err)
	assert.NotNil(t, stats)
	assert.False(t, stats.IsDir())
}

func TestNewFile_absolutePathFileDoesntExists(t *testing.T) {
	name := "/tmp/someveryweirdnamethatshouldntexistpls"
	file, err := NewFile(name)
	assert.Nil(t, err)

	assert.False(t, file.Exists())
	assert.Equal(t, "someveryweirdnamethatshouldntexistpls", file.Name())
	assert.Equal(t, name, file.Path())
	stats, err := file.Stats()
	assert.Nil(t, stats)
	assert.Errorf(t, err, "file: file %s doesn't exists", name)
}

func TestFile_parent(t *testing.T) {
	file, err := NewFile("/tmp/lala")
	assert.Nil(t, err)

	parent1 := file.Parent()
	assert.Equal(t, "/tmp", parent1.Path())

	parent2 := parent1.Parent()
	assert.Equal(t, "/", parent2.Path())

	parent3 := parent2.Parent()
	assert.Nil(t, parent3)
}
