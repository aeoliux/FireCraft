package main

import (
	"log"
	"os"

	"github.com/therecipe/qt/core"
	"github.com/therecipe/qt/widgets"
	"github.com/zapomnij/firecraft/bin/firecraft/window"
	"github.com/zapomnij/firecraft/pkg/downloader"
)

var qApp *widgets.QApplication

func main() {
	downloader.MakeAllDirs(nil)
	if err := os.Chdir(downloader.LauncherDir); err != nil {
		log.Println("WARN: failed to change to .minecraft/launcher directory")
	}

	qApp = widgets.NewQApplication(len(os.Args), os.Args)
	core.QCoreApplication_SetApplicationName("FireCraft")
	core.QCoreApplication_SetApplicationVersion("1.0.0")

	var fw = window.NewFWindow()
	fw.Window.Resize2(800, 600)

	fw.Window.Show()
	qApp.Exec()
}
