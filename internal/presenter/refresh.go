package presenter

import "github.com/google/logger"

type viewUpdateEvent struct {
	data viewData
	cfg  viewConfig
}

func refresh(ev viewUpdateEvent) {
	go func() {
		logger.Infof("Running")
		if ev.data.toRefresh&LocationBar != 0 {
			ev.data.view.SetLocationBar(ev.data.location)
		}

		if ev.data.toRefresh&FileListView != 0 {
			ev.data.view.SetFileList(ev.data.fileList)
		}

		if ev.data.toRefresh&ContextView != 0 {
			if ev.cfg.showSelection {
				ev.data.view.SetContextDetails(ev.data.selectedList)
			} else {
				if ev.data.fileDetail != nil {
					ev.data.view.SetContextDetails(*ev.data.fileDetail)
				} else {
					ev.data.view.SetContextDetails(nil)
				}
			}
		}

		if ev.data.toRefresh&StatusBar != 0 {
			ev.data.view.SetStatusMessage(ev.data.status)
		}

		if ev.data.toRefresh&CommandBar != 0 {
			ev.data.view.SetCommandBar(ev.data.commandBarContent)
		}

		if ev.data.toRefresh&Focus != 0 {
			if ev.cfg.mode == Navigation {
				ev.data.view.FocusFileList()
			} else {
				ev.data.view.FocusCommandBar()
			}
		}
	}()
}
