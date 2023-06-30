package downloader

import (
	"fmt"
	"os"
	"path"
	"strings"

	"github.com/zapomnij/firecraft/pkg/unzip"
)

func (v VersionJSON) FetchLibraries(log chan string) error {
	verClass := path.Join(VersionDir, v.Id, v.Id+".jar")
	classpath := ""

	if _, err := DownloadFile(v.Downloads.Client.Url, verClass); err != nil {
		log <- "error: " + err.Error()
		return err
	}

	if err := os.MkdirAll(path.Join(NativesDir, v.Id), os.ModePerm); err != nil {
		log <- "error: " + err.Error()
		return err
	}

	log <- "downloader: analyzing libraries. This can take a while"
	for _, j := range v.Libraries {
		if j.Rules != nil {
			rules := *(j.Rules)
			passed := false
			for _, rule := range rules {
				if rule.Action == "allow" {
					if (rule.Os != nil && rule.Os.Name == OperatingSystem) || rule.Os == nil {
						passed = true
					}
				} else {
					if rule.Os.Name == OperatingSystem {
						passed = false
					}
				}
			}

			if !passed {
				continue
			}
		}
		if j.Downloads.Artifact != nil {
			pth := path.Join(LibrariesDir, j.Downloads.Artifact.Path)
			if !IsExist(pth) {
				log <- "downloader: downloading library '" + *j.Name + "'"
				if err := getLib(pth, j.Downloads.Artifact.Url); err != nil {
					log <- "error: " + err.Error()
					return err
				}
			}
			classpath += ":" + pth
		}

		if j.Downloads.Artifact == nil && j.Natives == nil && j.Name != nil {
			log <- "downloader: adding local library '" + *j.Name + "'"

			split := strings.Split(*j.Name, ":")
			category := split[0]
			name := split[1]
			ver := split[2]

			classpath += ":" + path.Join(LibrariesDir, category, name, ver, fmt.Sprintf("%s-%s.jar", name, ver))
		}

		if j.Natives != nil {
			if native, ok := (*j.Natives)[OperatingSystem]; ok {
				classifier := (*j.Downloads.Classifiers)[native]
				pth := path.Join(LibrariesDir, classifier.Path)
				if !IsExist(pth) {
					log <- "downloader: downloading native '" + classifier.Path + "'"
					if err := getLib(pth, classifier.Url); err != nil {
						log <- "error: " + err.Error()
						return err
					}
				}
				classpath += ":" + pth

				if j.Extract != nil {
					exclude := []string{}
					if j.Extract.Exclude != nil {
						exclude = *(j.Extract.Exclude)
					}

					nativesDir := path.Join(NativesDir, v.Id)
					log <- "downloader: extracting native '" + classifier.Path + "'"
					if err := unzip.Unzip(pth, nativesDir, exclude); err != nil {
						log <- "error: " + err.Error()
						return err
					}
				}
			}
		}
	}

	classpath += ":" + verClass
	log <- classpath
	return nil
}

func getLib(out, url string) error {
	outDir := path.Dir(out)
	if err := os.MkdirAll(outDir, os.ModePerm); err != nil {
		return err
	}

	_, err := DownloadFile(url, out)
	return err
}
