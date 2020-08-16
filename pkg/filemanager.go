package pkg

import (
	"fmt"
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

func (fm *FileManager) GetContextIndex() int {
	return fm.currentContext
}

func (fm *FileManager) GetContextsCount() int {
	return len(fm.contexts)
}

func (fm *FileManager) NewContext(baseDir string) error {
	ctx, err := createFileContext(baseDir)
	if err != nil {
		return err
	} else {
		fm.contexts = append(fm.contexts, ctx)
		return nil
	}
}

func (fm *FileManager) SwitchContext(index int) error {
	if index < 0 || index >= len(fm.contexts) {
		return fmt.Errorf("fm: context out of bounds: %d", index)
	}
	fm.currentContext = index
	return nil
}

func (fm *FileManager) DiscardContext(index int) error {
	if index < 0 || index >= len(fm.contexts) {
		return fmt.Errorf("fm: context out of bounds: %d", index)
	} else if fm.GetContextsCount() == 1 {
		return fmt.Errorf("fm: only one context exists, unable to close")
	}
	for i := index; i < len(fm.contexts)-1; i++ {
		fm.contexts[i] = fm.contexts[i+1]
	}
	fm.contexts = fm.contexts[:len(fm.contexts)-1]
	if fm.currentContext >= len(fm.contexts) {
		fm.currentContext = fm.currentContext - 1
	}
	return nil
}

func (fm *FileManager) GetCWD() vfs.File {
	return fm.context().GetCWD()
}

func (fm *FileManager) ListFile() ([]vfs.File, error) {
	return fm.context().ListFiles()
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

func createFileContext(baseDir string) (*filecontext.FileContext, error) {
	fs, err := localfs.NewFS(baseDir) // TODO create different fs depending on the basedir
	if err != nil {
		return nil, err
	}
	return filecontext.NewFileContext(fs), nil
}

func NewFileManager(baseDir string) *FileManager {
	ctx, err := createFileContext(baseDir)
	if err != nil {
		log.Fatalf("Error initializing FileManager: %v", err)
	}
	return &FileManager{
		contexts:       []*filecontext.FileContext{ctx},
		currentContext: 0,
		selection:      vfs.NewFileSet(),
	}
}
