package downloader

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/zapomnij/firecraft/pkg/unzip"
)

func (v VersionJSON) FetchLibraries(log chan string) error {
	var classpathSeparator string
	if OperatingSystem == "windows" {
		classpathSeparator = ";"
	} else {
		classpathSeparator = ":"
	}

	verClass := filepath.Join(VersionDir, v.Id, v.Id+".jar")
	classpath := ""

	if v.Downloads != nil && !CheckSumByPath(verClass, v.Downloads.Client.Sha1) {
		log <- "downloader: " + v.Id + ".jar"
		if err := DownloadAndCheck(verClass, v.Downloads.Client.Url, v.Downloads.Client.Sha1); err != nil {
			log <- "error: " + err.Error()
			return err
		}
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
					if rule.Os == nil || rule.Os.Name == OperatingSystem {
						passed = false
					}
				}
			}

			if !passed {
				continue
			}
		}
		if j.Downloads != nil && j.Downloads.Artifact != nil {
			pth := filepath.Join(LibrariesDir, j.Downloads.Artifact.Path)
			if !CheckSumByPath(pth, j.Downloads.Artifact.Sha1) {
				log <- "downloader: downloading library '" + *j.Name + "'"
				if err := DownloadAndCheck(path.Join(strings.Split(pth, "\\")...), j.Downloads.Artifact.Url, j.Downloads.Artifact.Sha1); err != nil {
					log <- "error: " + err.Error()
					return err
				}
			}
		}

		if j.Natives == nil && j.Name != nil {
			log <- "downloader: adding library '" + *j.Name + "'"

			split := strings.Split(*j.Name, ":")
			category := []string{LibrariesDir}
			category = append(category, strings.Split(split[0], ".")...)
			name := split[1]
			ver := split[2]

			classpath += classpathSeparator + filepath.Join(filepath.Join(category...), name, ver, fmt.Sprintf("%s-%s.jar", name, ver))
		}

		if j.Natives != nil {
			if native, ok := (*j.Natives)[OperatingSystem]; ok {
				var pth string
				if j.Downloads != nil {
					classifier := (*j.Downloads.Classifiers)[native]
					pth = filepath.Join(LibrariesDir, classifier.Path)
					if !CheckSumByPath(pth, classifier.Sha1) {
						log <- "downloader: downloading native '" + classifier.Path + "'"
						if err := DownloadAndCheck(path.Join(strings.Split(pth, "\\")...), classifier.Url, classifier.Sha1); err != nil {
							log <- "error: " + err.Error()
							return err
						}
					}
				} else if j.Name != nil {
					split := strings.Split(*j.Name, ":")
					name := split[1]
					ver := split[2]
					category := []string{LibrariesDir}
					category = append(category, strings.Split(split[0], ".")...)

					pth = filepath.Join(filepath.Join(category...), name, ver, fmt.Sprintf("%s-%s-%s.jar", name, ver, native))
				} else {
					continue
				}

				classpath += classpathSeparator + pth

				if j.Extract != nil {
					exclude := []string{}
					if j.Extract.Exclude != nil {
						exclude = *(j.Extract.Exclude)
					}

					log <- "downloader: extracting native '" + pth + "'"
					if err := unzip.Unzip(pth, NativesDir, exclude); err != nil {
						log <- "error: " + err.Error()
						return err
					}
				}
			}
		}
	}

	classpath += classpathSeparator + verClass
	log <- classpath

	return nil
}
