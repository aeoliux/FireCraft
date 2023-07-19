package downloader

import (
	"os"
	"path/filepath"
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
		homedir = os.Getenv("HOME")
		mcDir := path.Join(homedir, "Library", "Application Support", "minecraft")
		dir = &mcDir
	}

	if dir != nil {
		MinecraftDir = *dir
	} else if homedir == "" {
		MinecraftDir = filepath.Join(".", ".minecraft")
	} else {
		MinecraftDir = filepath.Join(homedir, ".minecraft")
	}

	VersionDir = filepath.Join(MinecraftDir, "versions")
	AssetsDir = filepath.Join(MinecraftDir, "assets")
	LibrariesDir = filepath.Join(MinecraftDir, "libraries")
	NativesDir = filepath.Join(MinecraftDir, "natives")
	LauncherDir = filepath.Join(MinecraftDir, "launcher")
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
