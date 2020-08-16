package presenter

import (
	"github.com/daniloqueiroz/dfm/internal/view"
	"github.com/daniloqueiroz/dfm/pkg"
	"github.com/daniloqueiroz/dfm/pkg/vfs"
	"github.com/google/logger"
	"log"
	"strconv"
	"strings"
)

type AppMode string

type presenter struct {
	quitFunc  func()
	fm        *pkg.FileManager
	cfg       *viewConfig
	data      *viewData
	refresher *viewRefresher
}

func (p *presenter) Start() {
	cwd := p.fm.GetCWD()
	p.data.location = cwd.Path()
	p.data.fileList = p.getFiles(cwd)
	p.data.status = p.getStatus()
	p.data.fileDetail = nil
	p.data.selectedList = p.getSelectedItems()
	p.data.commandBarContent = ""
	p.refresher.start()
}

func (p *presenter) onEvent(event interface{}) {
	logger.Infof("event received: %T", event)
	switch ev := event.(type) {
	case view.ScreenSizeChanged:
		p.data.toRefresh = ALL
	case view.Command:
		p.onCommandInput(ev.Cmdline)
	case view.ToggleCommandMode:
		p.onCommandModeEvent(ev)
	case view.NavPrev:
		p.onNavPrevEvent()
	case view.NavNext:
		p.onNavNextEvent()
	case view.FileListItemHover:
		p.onListListHoverEvent(ev)
	case view.FileListItemSelected:
		p.onItemSelectedEvent(ev)
	case view.SwitchContext:
		p.onChangeContext(ev.Index)
	default:
		logger.Infof("Unhandled event: %v", ev)
	}
	p.refresher.refresh(viewRefreshEvent{
		data: *p.data,
		cfg:  *p.cfg,
	})
}

func (p *presenter) onCommandInput(cmdline string) {
	tokens := strings.Split(cmdline, " ")
	cmd := tokens[0]
	params := strings.TrimSpace(cmdline[len(cmd):])
	switch cmd {
	case QuitCommand, QuitCommandShort:
		p.quitFunc()
	case NewContextShort:
		p.onNewContext("")
	case CloseContextShort:
		p.onCloseContext(p.fm.GetContextIndex())
	case SelectionToggleCommandShort:
		p.onToggleSelectionViewEvent()
	case HiddenToggleCommandShort:
		p.onFileVisibilityEvent()
	case ToggleCommand:
		if params == "selection" || params == "details" {
			p.onToggleSelectionViewEvent()
		} else if params == "hidden" {
			p.onFileVisibilityEvent()
		}
	case Context:
		if len(tokens) < 2 {
			return
		} else {
			subcommand := tokens[1]
			switch subcommand {
			case "new":
				if len(tokens) > 2 {
					params = tokens[2]
				} else {
					params = ""
				}
				p.onNewContext(params)
			case "close":
				var idx int
				if len(tokens) > 2 {
					var err error
					idx, err = strconv.Atoi(tokens[2])
					if err != nil {
						return
					}
				} else {
					idx = p.fm.GetContextIndex()
				}
				p.onCloseContext(idx)
			}
		}
	}
	p.data.toRefresh = p.data.toRefresh | Focus | CommandBar
	p.cfg.mode = Navigation
}

func (p *presenter) onNavNextEvent() {
	p.fm.NavNext()
	cwd := p.fm.GetCWD()
	p.data.location = cwd.Path()
	p.data.fileList = p.getFiles(cwd)
	p.data.selectedList = p.getSelectedItems()
	p.data.status = p.getStatus()
	p.data.fileDetail = nil
	p.data.toRefresh = ALL
}

func (p *presenter) onNavPrevEvent() {
	p.fm.NavPrev()
	cwd := p.fm.GetCWD()
	p.data.location = cwd.Path()
	p.data.fileList = p.getFiles(cwd)
	p.data.selectedList = p.getSelectedItems()
	p.data.status = p.getStatus()
	p.data.fileDetail = nil
	p.data.toRefresh = ALL
}

func (p *presenter) onFileVisibilityEvent() {
	p.cfg.hideHidden = !p.cfg.hideHidden
	cwd := p.fm.GetCWD()
	p.data.fileList = p.getFiles(cwd)
	p.data.status = p.getStatus()
	p.data.toRefresh = FileListView | StatusBar
}

func (p *presenter) onToggleSelectionViewEvent() {
	p.cfg.showSelection = !p.cfg.showSelection
	p.data.toRefresh = ContextView
}

func (p *presenter) onListListHoverEvent(ev view.FileListItemHover) {
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
}

func (p *presenter) onCommandModeEvent(ev view.ToggleCommandMode) {
	p.data.toRefresh = Focus | CommandBar
	p.data.commandBarContent = ""
	if p.cfg.mode != Navigation {
		p.cfg.mode = Navigation
	} else {
		p.cfg.mode = Command
		if ev.Prefix != 0 {
			p.data.commandBarContent = string(ev.Prefix)
		}
	}
}

func (p *presenter) onItemSelectedEvent(ev view.FileListItemSelected) {
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
		ContextCount:       p.fm.GetContextsCount(),
		ActiveContext:      p.fm.GetContextIndex(),
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

func (p *presenter) onNewContext(baseDir string) {
	if baseDir == "" {
		baseDir = p.fm.GetCWD().Path()
	}
	err := p.fm.NewContext(baseDir) // TODO handle error
	if err == nil {
		p.onChangeContext(p.fm.GetContextsCount() - 1)
	}
}

func (p *presenter) onChangeContext(idx int) {
	err := p.fm.SwitchContext(idx) // TODO handle error
	if err == nil {
		cwd := p.fm.GetCWD()
		p.data.location = cwd.Path()
		p.data.fileList = p.getFiles(cwd)
		p.data.selectedList = p.getSelectedItems()
		p.data.status = p.getStatus()
		p.data.fileDetail = nil
		p.data.toRefresh = ALL
	}
}

func (p *presenter) onCloseContext(idx int) {
	err := p.fm.DiscardContext(idx) // TODO handle error
	if err == nil {
		cwd := p.fm.GetCWD()
		p.data.location = cwd.Path()
		p.data.fileList = p.getFiles(cwd)
		p.data.selectedList = p.getSelectedItems()
		p.data.status = p.getStatus()
		p.data.fileDetail = nil
		p.data.toRefresh = ALL
	}
}

func NewPresenter(fm *pkg.FileManager, view view.View) *presenter {
	p := &presenter{
		quitFunc: view.Quit,
		fm:       fm,
		cfg: &viewConfig{
			mode:          Navigation,
			hideHidden:    true,
			showSelection: false,
		},
		data:      &viewData{},
		refresher: newViewRefresher(view),
	}
	view.OnEvent(p.onEvent)

	return p
}
