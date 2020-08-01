package internal

import (
	"github.com/daniloqueiroz/dfm/internal/view"
	"github.com/daniloqueiroz/dfm/pkg"
	"github.com/daniloqueiroz/dfm/pkg/vfs"
	"github.com/google/logger"
	"log"
)

type config struct {
	hideHidden bool
}

type presenter struct {
	v   view.View
	fm  *pkg.FileManager
	cfg *config
}

func (p presenter) onEvent(event interface{}) {
	msg := ""
	refresh := false
	switch ev := event.(type) {
	case view.ViewSizeChanged:
		refresh = true
	case view.NavPrev:
		p.fm.NavPrev()
		refresh = true
	case view.NavNext:
		p.fm.NavNext()
		refresh = true
	case view.FileListItemSelected:
		// try cd file
		err := p.fm.CD(ev.Name)
		if err != nil {
			// file isn't a dir or doesnt exists
			logger.Warningf("Can't change dir: %v", err)
			msg = "Can't change dir" // TODO clear msg after x secs
		}
		refresh = true
	default:
		logger.Infof("Unhandled event: %v", ev)
	}
	if refresh {
		p.refresh(msg)
	}
}

func (p presenter) refresh(msg string) {
	cwd := p.fm.GetCWD()
	p.v.SetCurrentDir(cwd.Path())
	items := p.getFiles(cwd)
	p.v.UpdateFileList(items)
	p.v.SetStatusMessage(view.Status{
		Context:            p.fm.GetContextNumber(),
		Message:            msg,
		FilesCount:         len(items),
		SelectedFilesCount: p.fm.CountSelectedFiles(),
	})
}

func (p presenter) getFiles(cwd vfs.File) []view.FileItem {
	files, err := p.fm.ListFile()
	if err != nil {
		log.Fatalf("Error listing files: %v", err)
	}
	var items []view.FileItem
	for _, file := range files {
		stats, err := file.Stats()
		if err != nil {
			log.Fatalf("Error getting file stats: %v", err)
		}
		if !p.cfg.hideHidden || !file.IsHidden() {
			items = append(items, view.FileItem{
				Name:     file.Name(),
				FullPath: file.Path(),
				IsDir:    stats.IsDir(),
				Size:     stats.Size(),
			})
		}
	}
	parent := cwd.Parent()
	if parent != nil {
		stats, _ := parent.Stats()
		items = append(items, view.FileItem{
			Name:     "..",
			FullPath: parent.Path(),
			IsDir:    stats.IsDir(),
			Size:     stats.Size(),
		})
	}
	return items
}

func (p presenter) Start() {
	p.v.Show(p.onEvent)
}

func NewPresenter(fm *pkg.FileManager, view view.View) *presenter {
	return &presenter{
		v:  view,
		fm: fm,
		cfg: &config{
			hideHidden: true,
		},
	}
}
