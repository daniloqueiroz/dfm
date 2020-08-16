package presenter

import (
	"github.com/daniloqueiroz/dfm/internal/view"
	"github.com/google/logger"
)

type viewStateChanged struct {
	data viewData
	cfg  viewConfig
}

type ViewDispatcher struct {
	StateChangedCh chan viewStateChanged
	view           view.View
}

func (vd *ViewDispatcher) Start() {
	go func() {
		for ev := range vd.StateChangedCh {
			logger.Infof("Refreshing view")
			data := ev.data
			cfg := ev.cfg

			if data.toRefresh&LocationBar != 0 {
				vd.view.SetLocationBar(data.location)
			}

			if data.toRefresh&FileListView != 0 {
				vd.view.SetFileList(data.fileList)
			}

			if data.toRefresh&ContextView != 0 {
				if cfg.showSelection {
					vd.view.SetContextDetails(data.selectedList)
				} else {
					if data.fileDetail != nil {
						vd.view.SetContextDetails(*data.fileDetail)
					} else {
						vd.view.SetContextDetails(nil)
					}
				}
			}

			if data.toRefresh&StatusBar != 0 {
				vd.view.SetStatusMessage(data.status)
			}

			if data.toRefresh&CommandBar != 0 {
				vd.view.SetCommandBar(data.commandBarContent)
			}

			if data.toRefresh&Focus != 0 {
				if cfg.mode == Navigation {
					vd.view.FocusFileList()
				} else {
					vd.view.FocusCommandBar()
				}
			}
		}
	}()
}

func NewViewDispatcher(view view.View) *ViewDispatcher {
	return &ViewDispatcher{
		view:           view,
		StateChangedCh: make(chan viewStateChanged, 2),
	}
}
