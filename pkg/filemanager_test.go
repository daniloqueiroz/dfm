package pkg

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func fileManager() *FileManager {
	return NewFileManager("/tmp")
}

func TestNewFileManager_CreateContext(t *testing.T) {
	/* TODO context tests
	 * Create new context (base folder and no base folder (use cwd from current ctx)
	 * Switch active context (corner cases-> out of bounds)
	 * Get number of contexts
	 * close/discard context
	 */
	fm := fileManager()
	assert.Equal(t, 1, fm.GetContextsCount())

	fm.NewContext("/proc")
	assert.Equal(t, 2, fm.GetContextsCount())
}

func TestNewFileManager_SwitchContext(t *testing.T) {
	fm := fileManager()
	fm.NewContext("/proc")
	assert.Equal(t, 0, fm.GetContextIndex())

	fm.SwitchContext(1)
	assert.Equal(t, 1, fm.GetContextIndex())
}

func TestNewFileManager_SwitchContext_outOfBounds(t *testing.T) {
	fm := fileManager()
	assert.Equal(t, 0, fm.GetContextIndex())
	err := fm.SwitchContext(1)
	assert.Error(t, err)
	err = fm.SwitchContext(-1)
	assert.Error(t, err)
}

func TestNewFileManager_SwitchContextChangesCWD(t *testing.T) {
	fm := fileManager()
	fm.NewContext("/proc")
	fm.SwitchContext(1)

	assert.Equal(t, "/proc", fm.GetCWD().Path())
}

func TestNewFileManager_CloseContext(t *testing.T) {
	fm := fileManager()
	fm.NewContext("/proc")
	assert.Equal(t, 2, fm.GetContextsCount())
	fm.NewContext("/")
	assert.Equal(t, 3, fm.GetContextsCount())

	fm.DiscardContext(1)
	assert.Equal(t, 2, fm.GetContextsCount())
}

func TestNewFileManager_CloseLastContext(t *testing.T) {
	fm := fileManager()
	fm.NewContext("/proc")
	assert.Equal(t, 2, fm.GetContextsCount())
	fm.NewContext("/")
	assert.Equal(t, 3, fm.GetContextsCount())
	fm.SwitchContext(2)

	fm.DiscardContext(2)
	assert.Equal(t, 2, fm.GetContextsCount())
	assert.Equal(t, 1, fm.GetContextIndex())
}

func TestNewFileManager_CloseCurrentContextFirst(t *testing.T) {
	fm := fileManager()
	fm.NewContext("/proc")
	assert.Equal(t, 2, fm.GetContextsCount())
	assert.Equal(t, 0, fm.GetContextIndex())

	fm.DiscardContext(0)
	assert.Equal(t, 1, fm.GetContextsCount())
	assert.Equal(t, 0, fm.GetContextIndex())
}

func TestNewFileManager_CloseCurrentContextOnlyOneContext(t *testing.T) {
	fm := fileManager()
	assert.Equal(t, 0, fm.GetContextIndex())

	err := fm.DiscardContext(0)
	assert.Error(t, err)

}
