package downloader

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"path"
)

type VersionInfo struct {
	Id              string
	Type            string
	Url             string
	Time            string
	ReleaseTime     string
	Sha1            string
	ComplianceLevel uint
}

type VersionManifest struct {
	Latest struct {
		Release  string
		Snapshot string
	}
	Versions []VersionInfo
}

func GetVersionManifest() (*VersionManifest, error) {
	out := path.Join(VersionDir, "version_manifest_v2.json")
	//Log <- "downloader: downloading version manifest"
	img, err := DownloadFile("https://piston-meta.mojang.com/mc/game/version_manifest_v2.json", out)
	if err != nil {
		if len(img) == 0 {
			return nil, err
		} else {
			log.Println("Failed to save version_manifest_v2.json")
		}
	}

	var vm VersionManifest
	err = json.Unmarshal(img, &vm)
	vm.appendLocalVersions()
	return &vm, err
}

func (vm *VersionManifest) appendLocalVersions() {
	dirs, err := ioutil.ReadDir(VersionDir)
	if err != nil {
		return
	}

l1:
	for _, dir := range dirs {
		if !dir.IsDir() || dir.Name() == "version_manifest_v2.json" {
			continue
		}
		for _, j := range vm.Versions {
			if j.Id == dir.Name() {
				continue l1
			}
		}

		stat, err := os.Stat(path.Join(VersionDir, dir.Name(), dir.Name()+".json"))
		if os.IsNotExist(err) || stat.IsDir() {
			continue
		}

		new := VersionInfo{
			Id:   dir.Name(),
			Type: "local",
		}

		vm.Versions = append(vm.Versions, new)
	}
}
