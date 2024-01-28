package window

import (
	"fmt"
	"github.com/therecipe/qt/gui"
	"os"
	"path"
	"strings"

	"github.com/zapomnij/firecraft/pkg/downloader"
	"github.com/zapomnij/firecraft/pkg/javafind"
	"github.com/zapomnij/firecraft/pkg/runner"
)

func (fw *FWindow) Launch() {
	defer fw.end()
	fw.gameLogger.SetPlainText("")
	fw.ProgressBar.SetVisible(true)

	fw.updateProgressBar(0, "launcher: preparing environment")
	if err := os.Chdir(downloader.MinecraftDir); err != nil {
		fw.gameLogger.InsertPlainText("launcher: " + err.Error() + "\n")
		return
	}

	if err := os.MkdirAll(downloader.NativesDir, os.ModePerm); err != nil {
		fw.gameLogger.InsertPlainText("launcher: " + err.Error() + "\n")
		return
	}

	prof := LProfile{}
	selected := fw.profilesSelector.CurrentText()
	if selected == "New profile" {
		prof.JavaArgs = "-Xmx2048M"
		prof.LastVersionId = vm.Latest.Release
	} else {
		prof = lpf.Profiles[selected]
		if prof.GameDir != "" {
			_ = downloader.MakeAllDirs(&prof.GameDir)
		}
	}

	if err := os.RemoveAll(downloader.NativesDir); err != nil && !os.IsNotExist(err) {
		fw.gameLogger.InsertPlainText("launcher: couldn't delete natives directory. Game may not work properly\n")
	}

	accToken, uuid, haveBoughtTheGame := "", "", false
	if fw.usernameTv.Text() == "" {
		fw.gameLogger.InsertPlainText("launcher: missing username or game not authenticated with Microsoft account\n")
		return
	}

	if fw.ms.Authed {
		accToken = fw.ms.AccessToken
		uuid = fw.ms.Uuid
		haveBoughtTheGame = fw.ms.HaveBoughtTheGame
	}

	if wd, err := os.Getwd(); err == nil {
		fw.gameLogger.InsertPlainText("launcher: gameDir = " + wd + "\n")
	}

	fw.updateProgressBar(10, "launcher: getting version manifest")
	vJson, err := downloader.NewClientJSON(*vm, prof.LastVersionId)
	if err != nil {
		fw.gameLogger.InsertPlainText("launcher: " + err.Error() + "\n")
		return
	}

	fw.updateProgressBar(10, "launcher: finding Java")
	if prof.JavaDir == "" {
		jvDir := javafind.FindJava(vJson.JavaVersion.MajorVersion)
		if jvDir == nil {
			fw.gameLogger.InsertPlainText(fmt.Sprintf("launcher: please install Java version %d\n", vJson.JavaVersion.MajorVersion))
			return
		}

		prof.JavaDir = *jvDir
	}

	fw.updateProgressBar(5, "launcher: downloading libraries")
	classpath, ch := "", make(chan string)
	go vJson.FetchLibraries(ch)
	for {
		msg := <-ch

		if strings.HasPrefix(msg, "error: ") {
			fw.gameLogger.InsertPlainText("launcher: process failed '" + msg + "'\n")
			return
		} else if !strings.HasPrefix(msg, "downloader: ") {
			classpath = msg
			break
		}

		fw.updateProgressBar(0, msg)
	}
	fw.updateProgressBar(25, "launcher: downloading assets")
	go vJson.GetAssets(ch)
	for {
		msg := <-ch

		if strings.HasPrefix(msg, "error: ") {
			fw.gameLogger.InsertPlainText("launcher: process failed '" + msg + "'\n")
			return
		} else if msg == "downloader: assets finished" {
			break
		}

		fw.updateProgressBar(0, msg)
	}

	fw.updateProgressBar(40, "launcher: setting up logger")
	outputChannel := make(chan string)
	go func() {
		for {
			msg := <-outputChannel
			if msg == "EOF" {
				break
			}
			fmt.Print(msg)

			//fw.gameLogger.MoveCursor(gui.QTextCursor__End, gui.QTextCursor__MoveAnchor)
			fw.gameLogger.InsertPlainText(msg)

			shrinker := fw.gameLogger.ToPlainText()
			if len(shrinker) > 10240 {
				fw.gameLogger.InsertPlainText(shrinker[len(shrinker)-10240:])
				//fw.gameLogger.MoveCursor(gui.QTextCursor__End, gui.QTextCursor__MoveAnchor)
			}
		}
	}()

	run := runner.NewRunner(fw.usernameTv.Text(), path.Join(prof.JavaDir, "bin", "java"), classpath, prof.JavaArgs, *vJson, downloader.Ai, outputChannel)
	if uuid != "" && accToken != "" {
		run.SetUpMicrosoft(uuid, accToken, haveBoughtTheGame)
	}

	fw.updateProgressBar(10, "launcher: starting Minecraft")
	fw.Window.SetVisible(false)
	if err := run.Run(); err != nil {
		fw.gameLogger.InsertPlainText("launcher: " + err.Error() + "\n")
	} else {
		fw.gameLogger.InsertPlainText("launcher: game exited with 0 exit status (no error)\n")
	}

	//fw.end()
}

func (fw *FWindow) end() {
	downloader.SetUpVariables(nil)
	if err := os.Chdir(downloader.LauncherDir); err != nil {
		fw.gameLogger.InsertPlainText("launcher: please restart launcher due to '" + err.Error() + "'\n")
		fw.Window.SetVisible(true)
		return
	}

	fw.gameLogger.MoveCursor(gui.QTextCursor__End, gui.QTextCursor__MoveAnchor)
	fw.playBt.SetEnabled(true)
	fw.ProgressBar.SetValue(0)
	fw.ProgressBar.SetVisible(false)
	fw.Window.SetVisible(true)
}
