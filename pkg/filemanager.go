package pkg

import (
	"github.com/daniloqueiroz/dfm/pkg/filecontext"
	"github.com/daniloqueiroz/dfm/pkg/vfs"
	"github.com/daniloqueiroz/dfm/pkg/vfs/localfs"
	"log"
)

type FileManager struct {
	contexts       []*filecontext.FileContext
	currentContext int
	selection      []*vfs.File
}

func (m FileManager) context() *filecontext.FileContext {
	return m.contexts[m.currentContext]
}

func (m FileManager) GetCWD() vfs.File {
	return m.context().GetCWD()
}

func (m FileManager) ListFile() ([]vfs.File, error) {
	return m.context().ListFiles()
}

func (m FileManager) GetContextNumber() int {
	return m.currentContext
}

func (m FileManager) CountSelectedFiles() int {
	return len(m.selection)
}

func (m FileManager) NavPrev() {
	m.context().Previous()
}

func (m FileManager) NavNext() {
	m.context().Next()
}

func (m FileManager) CD(name string) error {
	return m.context().CD(name)
}

func NewFileManager(basedir string) *FileManager {
	fs, err := localfs.NewFS(basedir) // TODO create different fs depending on the basedir
	if err != nil {
		log.Fatalf("Error initializing FileManager: %v", err)
	}
	ctx := filecontext.NewFileContext(fs)
	return &FileManager{
		contexts:       []*filecontext.FileContext{ctx},
		currentContext: 0,
		selection:      make([]*vfs.File, 0),
	}
}
