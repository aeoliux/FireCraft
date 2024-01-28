package window

import (
	"fmt"
	"os"
	"path"
	"strings"

	"github.com/zapomnij/firecraft/pkg/downloader"
	"github.com/zapomnij/firecraft/pkg/javafind"
	"github.com/zapomnij/firecraft/pkg/runner"
)

func (fw *FWindow) Launch() {
	defer fw.end()

	fw.logger.Clear()

	if err := os.Chdir(downloader.MinecraftDir); err != nil {
		fw.appendToLog("launcher: failed to change to .minecraft directory\n")
		//fw.end()
		return
	}

	if err := os.MkdirAll(downloader.NativesDir, os.ModePerm); err != nil {
		fw.appendToLog("launcher: failed to create natives directory\n")
		return
	}

	prof := LProfile{}
	selected := fw.profilesSelector.CurrentText()
	if selected == "New profile" {
		fw.appendToLog("launcher: template profile selected. Launching latest version\n")
		prof.JavaArgs = "-Xmx2048M"
		prof.LastVersionId = vm.Latest.Release
	} else {
		prof = lpf.Profiles[selected]
		if prof.GameDir != "" {
			_ = downloader.MakeAllDirs(&prof.GameDir)
		}
	}

	fw.appendToLog("launcher: clearing natives\n")
	if err := os.RemoveAll(downloader.NativesDir); err != nil && !os.IsNotExist(err) {
		fw.appendToLog("launcher: failed to clear natives, game may not work properly\n")
	}

	accToken, uuid, haveBoughtTheGame := "", "", false
	if fw.usernameTv.Text() == "" {
		fw.appendToLog("launcher: username not specified\n")
		return
	}

	if fw.ms.Authed {
		accToken = fw.ms.AccessToken
		uuid = fw.ms.Uuid
		haveBoughtTheGame = fw.ms.HaveBoughtTheGame
	}

	if wd, err := os.Getwd(); err == nil {
		fw.appendToLog(fmt.Sprintf("launcher: current working directory '%s'. Selected version is %s\n", wd, prof.LastVersionId))
	}

	fw.appendToLog(fmt.Sprintf("launcher: downloading %s.json\n", prof.LastVersionId))
	vJson, err := downloader.NewClientJSON(*vm, prof.LastVersionId)
	if err != nil {
		fw.appendToLog(fmt.Sprintf("launcher: version %s error '%s'\n", prof.LastVersionId, err))
		//fw.end()
		return
	}

	if prof.JavaDir == "" {
		jvDir := javafind.FindJava(vJson.JavaVersion.MajorVersion)
		if jvDir == nil {
			fw.appendToLog("launcher: failed to find Java automatically. Specify Java binary path in profile\n")
			//fw.end()
			return
		}

		prof.JavaDir = *jvDir
	}

	classpath, assetsDone, ch := "", false, make(chan string)
	fw.appendToLog("launcher: starting downloader\n")
	go vJson.FetchLibraries(ch)
	go vJson.GetAssets(ch)
	for {
		msg := <-ch

		if strings.HasPrefix(msg, "error: ") {
			fw.appendToLog(fmt.Sprintf("launcher: process failed '%s'\n", msg))
			return
		} else if !strings.HasPrefix(msg, "downloader: ") {
			fw.appendToLog("downloader: libraries finished\n")
			classpath = msg
		} else if msg == "downloader: assets finished" {
			fw.appendToLog(msg + "\n")
			assetsDone = true
		} else {
			fw.appendToLog(msg + "\n")
		}

		if classpath != "" && assetsDone {
			break
		}
	}

	fw.appendToLog("launcher: starting Minecraft\n")
	run := runner.NewRunner(fw.usernameTv.Text(), path.Join(prof.JavaDir, "bin", "java"), classpath, prof.JavaArgs, *vJson, downloader.Ai)
	if uuid != "" && accToken != "" {
		run.SetUpMicrosoft(uuid, accToken, haveBoughtTheGame)
	}

	fw.Window.SetVisible(false)
	if err := run.Run(); err != nil {
		fw.appendToLog(fmt.Sprintf("launcher: %s\n", err))
	} else {
		fw.appendToLog("launcher: Minecraft exited with zero exit status (error hasn't occurred)\n")
	}

	//fw.end()
}

func (fw *FWindow) end() {
	downloader.SetUpVariables(nil)
	if err := os.Chdir(downloader.LauncherDir); err != nil {
		fw.appendToLog("launcher: failed to change to .minecraft/launcher directory\n")
	}

	fw.playBt.SetEnabled(true)
	fw.Window.SetVisible(true)
}
