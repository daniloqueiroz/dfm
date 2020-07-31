package filecontext

import (
	"github.com/daniloqueiroz/dfm/pkg/vfs/localfs"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNavigation_RegularFlow(t *testing.T) {
	initial, _ := localfs.NewFile("/tmp/initial")
	file1, _ := localfs.NewFile("/tmp/file1")
	file2, _ := localfs.NewFile("/tmp/file2")
	nav := NewNavigation(initial)

	nav.SetCurrent(file1)
	assert.Equal(t, file1, nav.Current())

	nav.SetCurrent(file2)
	assert.Equal(t, file2, nav.Current())

	nav.Previous()
	assert.Equal(t, file1, nav.Current())

	nav.Next()
	assert.Equal(t, file2, nav.Current())

	nav.Previous()
	assert.Equal(t, file1, nav.Current())

	nav.Previous()
	assert.Equal(t, initial, nav.Current())

	nav.Next()
	assert.Equal(t, file1, nav.Current())

	nav.Next()
	assert.Equal(t, file2, nav.Current())
}


func TestNavigation_SetCurrentClearNext(t *testing.T) {
	initial, _ := localfs.NewFile("/tmp/initial")
	file1, _ := localfs.NewFile("/tmp/file1")
	file2, _ := localfs.NewFile("/tmp/file2")
	nav := NewNavigation(initial)

	nav.SetCurrent(file1)
	assert.Equal(t, file1, nav.Current())

	nav.SetCurrent(file2)
	assert.Equal(t, file2, nav.Current())

	nav.Previous()
	assert.Equal(t, file1, nav.Current())

	nav.Previous()
	assert.Equal(t, initial, nav.Current())

	nav.SetCurrent(file2)
	assert.Equal(t, file2, nav.Current())

	nav.Previous()
	assert.Equal(t, initial, nav.Current())

	nav.Next()
	assert.Equal(t, file2, nav.Current())

	nav.Next()
	assert.Equal(t, file2, nav.Current())
}