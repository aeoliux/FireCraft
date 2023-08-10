package downloader

import (
	"encoding/json"
	"errors"
	"log"
	"os"
	"path"
)

type DownloadsEntry struct {
	Sha1 string
	Size uint
	Url  string
}

type LibraryDownloadEntry struct {
	Path string
	Sha1 string
	Size uint
	Url  string
}

type VersionJSON struct {
	Arguments struct {
		Game []interface{}
		Jvm  []interface{}
	} `json:"arguments"`
	AssetIndex *struct {
		Id        string
		Sha1      string
		Size      uint
		TotalSize uint
		Url       string
	} `json:"assetIndex"`
	Assets          string `json:"assets"`
	ComplianceLevel uint   `json:"complianceLeveL"`
	Downloads       *struct {
		Client          DownloadsEntry
		Client_mappings DownloadsEntry
		Server          DownloadsEntry
		Server_mappings DownloadsEntry
	} `json:"downloads"`
	Id           string  `json:"id"`
	InheritsFrom *string `json:"inheritsFrom,omitempty"`
	JavaVersion  struct {
		Component    string
		MajorVersion uint
	} `json:"javaVersion"`
	Libraries []struct {
		Downloads *struct {
			Artifact    *LibraryDownloadEntry
			Classifiers *map[string]struct {
				Path string `json:"path"`
				Size uint   `json:"size"`
				Sha1 string `json:"sha1"`
				Url  string `json:"url"`
			}
		} `json:"downloads"`
		Name    *string `json:"name"`
		Extract *struct {
			Exclude *[]string
		} `json:"extract"`
		Natives *map[string]string `json:"natives"`
		Rules   *[]struct {
			Action string
			Os     *struct {
				Name string
			}
		} `json:"rules"`
	} `json:"libraries"`
	Logging struct {
		Client struct {
			Argument string
			File     struct {
				Id   string
				Sha1 string
				Size uint
				Url  string
			}
			Type string
		}
	} `json:"logging"`
	MainClass              string      `json:"mainClass"`
	MinecraftArguments     *string     `json:"minecraftArguments,omitempty"`
	MinimumLauncherVersion interface{} `json:"minimumLauncherVersion"`
	ReleaseTime            string      `json:"releaseTime"`
	Time                   string      `json:"time"`
	Type                   string      `json:"type"`
}

func NewClientJSON(vm VersionManifest, version string) (*VersionJSON, error) {
	for _, j := range vm.Versions {
		if j.Id == version {
			var vj VersionJSON
			if j.Type == "local" {
				customVerPath := path.Join(VersionDir, version, version+".json")

				img, err := os.ReadFile(customVerPath)
				if err != nil {
					return nil, err
				}

				if err := json.Unmarshal(img, &vj); err != nil {
					return nil, err
				}

				if vj.InheritsFrom == nil {
					return &vj, nil
				}

				inheritsFrom := *vj.InheritsFrom
				ih, err := NewClientJSON(vm, inheritsFrom)
				if err != nil {
					return nil, err
				}

				mergeJSONs(ih, vj)
				return ih, nil
			} else {
				outdir := path.Join(VersionDir, version)
				if err := os.MkdirAll(outdir, os.ModePerm); err != nil {
					return nil, err
				}

				//log <- "downloader: downloading " + version + ".json"

				out := path.Join(outdir, version+".json")
				img, err := DownloadFile(j.Url, out)
				if err != nil {
					if len(img) == 0 {
						return nil, err
					} else {
						log.Printf("Failed to save %s.json: %s\n", version, err)
					}
				}

				if err = json.Unmarshal(img, &vj); err != nil {
					return nil, err
				}
			}

			return &vj, nil
		}
	}

	return nil, errors.New("Version not found")
}

func mergeJSONs(a *VersionJSON, b VersionJSON) {
	a.Libraries = append(a.Libraries, b.Libraries...)
	a.Arguments.Jvm = append(a.Arguments.Jvm, b.Arguments.Jvm...)
	a.Arguments.Game = append(a.Arguments.Game, b.Arguments.Game...)

	if a.Id != b.Id {
		a.Id = b.Id
	}

	if a.MainClass != b.MainClass {
		a.MainClass = b.MainClass
	}

	if a.MinecraftArguments != nil && b.MinecraftArguments != nil && a.MinecraftArguments != b.MinecraftArguments {
		a.MinecraftArguments = b.MinecraftArguments
	}

	if a.AssetIndex == nil {
		a.AssetIndex = b.AssetIndex
	}
}
