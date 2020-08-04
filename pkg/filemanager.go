package pkg

import (
	"github.com/daniloqueiroz/dfm/pkg/filecontext"
	"github.com/daniloqueiroz/dfm/pkg/vfs"
	"github.com/daniloqueiroz/dfm/pkg/vfs/localfs"
	"github.com/google/logger"
	"log"
	"os"
)

type FileManager struct {
	contexts       []*filecontext.FileContext
	currentContext int
	selection      vfs.FileSet
}

func (fm *FileManager) context() *filecontext.FileContext {
	return fm.contexts[fm.currentContext]
}

func (fm *FileManager) GetCWD() vfs.File {
	return fm.context().GetCWD()
}

func (fm *FileManager) ListFile() ([]vfs.File, error) {
	return fm.context().ListFiles()
}

func (fm *FileManager) GetContextNumber() int {
	return fm.currentContext
}

func (fm *FileManager) CountSelectedFiles() int {
	return fm.selection.Len()
}

func (fm *FileManager) NavPrev() {
	fm.context().Previous()
}

func (fm *FileManager) NavNext() {
	fm.context().Next()
}

func (fm *FileManager) CD(name string) error {
	return fm.context().CD(name)
}

func (fm *FileManager) Stats(name string) (os.FileInfo, error) {
	return fm.context().Stats(name)
}

func (fm *FileManager) SelectedFiles() []vfs.File {
	return fm.selection.Iterator()
}

func (fm *FileManager) Select(name string) {
	file, err := fm.context().GetFile(name)
	if err != nil {
		logger.Warningf("Unable to select file: %+v", err)
	} else {
		fm.selection.Add(file)
	}
}

func (fm *FileManager) Deselect(name string) {
	file, err := fm.context().GetFile(name)
	if err != nil {
		logger.Warningf("Unable to deselect file: %+v", err)
	} else {
		fm.selection.Remove(file)

	}
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
		selection:      vfs.NewFileSet(),
	}
}
