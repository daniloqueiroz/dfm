package presenter

import (
	"github.com/daniloqueiroz/dfm/internal/view"
	"github.com/daniloqueiroz/dfm/pkg"
	"github.com/daniloqueiroz/dfm/pkg/vfs"
	"github.com/google/logger"
	"log"
)

type config struct {
	hideHidden    bool
	showSelection bool
}

type presenter struct {
	v    view.View
	fm   *pkg.FileManager
	cfg  *config
	data *viewData
}

func (p *presenter) onEvent(event interface{}) {
	logger.Infof("ViewEvent received: %T", event)
	switch ev := event.(type) {
	case view.ScreenSizeChanged:
		p.data.toRefresh = ALL
	case view.Quit:
		p.v.Quit()
	case view.ToggleHiddenFilesVisibility:
		p.cfg.hideHidden = !p.cfg.hideHidden
		cwd := p.fm.GetCWD()
		p.data.fileList = p.getFiles(cwd)
		p.data.status = p.getStatus()
		p.data.toRefresh = FileListView | StatusBar
	case view.ToggleFileSelectionView:
		p.cfg.showSelection = !p.cfg.showSelection
		p.data.toRefresh = ContextView
	case view.NavPrev:
		p.fm.NavPrev()
		cwd := p.fm.GetCWD()
		p.data.location = cwd.Path()
		p.data.fileList = p.getFiles(cwd)
		p.data.selectedList = p.getSelectedItems()
		p.data.status = p.getStatus()
		p.data.fileDetail = nil
		p.data.toRefresh = ALL
	case view.NavNext:
		p.fm.NavNext()
		cwd := p.fm.GetCWD()
		p.data.location = cwd.Path()
		p.data.fileList = p.getFiles(cwd)
		p.data.selectedList = p.getSelectedItems()
		p.data.status = p.getStatus()
		p.data.fileDetail = nil
		p.data.toRefresh = ALL
	case view.FileListItemHover:
		stats, err := p.fm.Stats(ev.Name)
		if err != nil {
			logger.Warningf("Error retrieving %s stats", ev.Name)
		} else if stats != nil {
			p.data.fileDetail = &view.FileDetails{
				Size:    stats.Size(),
				Mode:    stats.Mode(),
				ModTime: stats.ModTime(),
			}
		}
		p.data.toRefresh = ContextView
	case view.FileListItemSelected:
		p.handleItemSelectedEvent(ev)
	default:
		logger.Infof("Unhandled event: %v", ev)
	}
	go refresh(p.v, *p.data, *p.cfg)
}

func (p *presenter) handleItemSelectedEvent(ev view.FileListItemSelected) {
	switch ev.SelectionMode {
	case view.Open:
		err := p.fm.CD(ev.Name)
		if err != nil {
			// file isn't a dir or doesnt exists
			logger.Warningf("Can't change dir: %v", err)
		}
		p.data.fileDetail = nil
		p.data.toRefresh = ALL
	case view.AddSelectionList:
		p.fm.Select(ev.Name)
		p.data.toRefresh = ContextView | StatusBar
	case view.RemoveSelectionList:
		p.fm.Deselect(ev.Name)
		p.data.toRefresh = ContextView | StatusBar
	}

	cwd := p.fm.GetCWD()
	p.data.location = cwd.Path()
	p.data.fileList = p.getFiles(cwd)
	p.data.selectedList = p.getSelectedItems()
	p.data.status = p.getStatus()
}

func (p *presenter) getStatus() view.Status {
	return view.Status{
		Context:            p.fm.GetContextNumber(),
		FilesCount:         len(p.data.fileList),
		SelectedFilesCount: len(p.data.selectedList),
	}
}

func (p *presenter) getSelectedItems() []view.FileItem {
	var items []view.FileItem
	for _, file := range p.fm.SelectedFiles() {
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

	return items
	// TODO update number of selected on status
}

func (p *presenter) getFiles(cwd vfs.File) []view.FileItem {
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

func (p *presenter) Start() {
	cwd := p.fm.GetCWD()
	p.data.location = cwd.Path()
	p.data.fileList = p.getFiles(cwd)
	p.data.status = p.getStatus()
	p.data.fileDetail = nil
	p.data.selectedList = p.getSelectedItems()
	p.v.Show(p.onEvent)
}

func NewPresenter(fm *pkg.FileManager, view view.View) *presenter {
	return &presenter{
		v:  view,
		fm: fm,
		cfg: &config{
			hideHidden:    true,
			showSelection: false,
		},
		data: &viewData{},
	}
}
