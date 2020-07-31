package internal

import (
	"github.com/google/logger"
)

func OnPanic(name string) {
	r := recover()
	if r != nil {
		logger.Errorf("Panic on %s: %+v", name, r)
	}
}
