package presenter

import (
	"github.com/daniloqueiroz/dfm/internal/view"
	"github.com/google/logger"
)

type viewRefreshEvent struct {
	data viewData
	cfg  viewConfig
}

type viewRefresher struct {
	evChan chan viewRefreshEvent
	view   view.View
}

func (vr *viewRefresher) start() {
	go func() {
		for ev := range vr.evChan {
			logger.Infof("Refreshing view")
			data := ev.data
			cfg := ev.cfg

			if data.toRefresh&LocationBar != 0 {
				vr.view.SetLocationBar(data.location)
			}

			if data.toRefresh&FileListView != 0 {
				vr.view.SetFileList(data.fileList)
			}

			if data.toRefresh&ContextView != 0 {
				if cfg.showSelection {
					vr.view.SetContextDetails(data.selectedList)
				} else {
					if data.fileDetail != nil {
						vr.view.SetContextDetails(*data.fileDetail)
					} else {
						vr.view.SetContextDetails(nil)
					}
				}
			}

			if data.toRefresh&StatusBar != 0 {
				vr.view.SetStatusMessage(data.status)
			}

			if data.toRefresh&CommandBar != 0 {
				vr.view.SetCommandBar(data.commandBarContent)
			}

			if data.toRefresh&Focus != 0 {
				if cfg.mode == Navigation {
					vr.view.FocusFileList()
				} else {
					vr.view.FocusCommandBar()
				}
			}
		}
	}()
	vr.view.Show()
}

func (vr *viewRefresher) refresh(ev viewRefreshEvent) {
	vr.evChan <- ev
}

func newViewRefresher(view view.View) *viewRefresher {
	return &viewRefresher{
		view:   view,
		evChan: make(chan viewRefreshEvent, 2),
	}
}
