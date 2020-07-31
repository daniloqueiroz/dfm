package main

import (
	"github.com/akamensky/argparse"
	"github.com/daniloqueiroz/dfm/internal"
	"github.com/daniloqueiroz/dfm/internal/cui"
	"github.com/daniloqueiroz/dfm/pkg"
	"github.com/google/logger"
	"io/ioutil"
	"os"
)

func main() {
	// load config
	logger.Init("dfm", false, true, ioutil.Discard)
	defer logger.Close()

	cwd, err := os.Getwd() // TODO receive it as parameter
	if err != nil {
		logger.Fatalf("Unable to get current directory", err)
	}
	parser := argparse.NewParser("dfm", "file manager")
	startDir := parser.String("d", "directory", &argparse.Options{
		Required: false,
		Help:     "start directory",
		Default:  cwd,
	})
	err = parser.Parse(os.Args)
	if err != nil {
		logger.Fatal("Error parsing parameters", err)
	}

	fm := pkg.NewFileManager(*startDir)
	w := cui.NewWindow()
	p := internal.NewPresenter(fm, w)
	logger.Info("Starting dfm")
	p.Start()
}
