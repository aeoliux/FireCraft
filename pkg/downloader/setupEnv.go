package downloader

import (
	"os"
	"path"
	"runtime"
)

var (
	OperatingSystem,
	MinecraftDir,
	VersionDir,
	LibrariesDir,
	NativesDir,
	LauncherDir,
	AssetsDir string
)

func SetUpVariables(dir *string) {
	var homedir string
	switch runtime.GOOS {
	case "linux", "windows":
		OperatingSystem = runtime.GOOS
		if runtime.GOOS == "windows" {
			homedir, _ = os.UserConfigDir()
		} else {
			homedir = os.Getenv("HOME")
		}
	case "darwin":
		OperatingSystem = "osx"
	}

	if dir != nil {
		MinecraftDir = *dir
	} else if homedir == "" {
		MinecraftDir = path.Join(".", ".minecraft")
	} else {
		MinecraftDir = path.Join(homedir, ".minecraft")
	}

	VersionDir = path.Join(MinecraftDir, "versions")
	AssetsDir = path.Join(MinecraftDir, "assets")
	LibrariesDir = path.Join(MinecraftDir, "libraries")
	NativesDir = path.Join(MinecraftDir, "natives")
	LauncherDir = path.Join(MinecraftDir, "launcher")
}

func MakeAllDirs(dir *string) error {
	SetUpVariables(dir)

	if err := os.MkdirAll(MinecraftDir, os.ModePerm); err != nil {
		return err
	}
	if err := os.MkdirAll(VersionDir, os.ModePerm); err != nil {
		return err
	}
	if err := os.MkdirAll(AssetsDir, os.ModePerm); err != nil {
		return err
	}
	if err := os.MkdirAll(LibrariesDir, os.ModePerm); err != nil {
		return err
	}
	if err := os.MkdirAll(NativesDir, os.ModePerm); err != nil {
		return err
	}
	if err := os.MkdirAll(LauncherDir, os.ModePerm); err != nil {
		return err
	}

	return nil
}
