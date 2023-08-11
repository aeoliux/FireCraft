package downloader

import (
	"crypto/sha1"
	"encoding/hex"
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

	if v.Downloads != nil && !checkSum(verClass, v.Downloads.Client.Sha1) {
		log <- "downloader: " + v.Id + ".jar"
		if _, err := DownloadFile(v.Downloads.Client.Url, verClass); err != nil {
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
			if !checkSum(pth, j.Downloads.Artifact.Sha1) {
				log <- "downloader: downloading library '" + *j.Name + "'"
				if err := getLib(path.Join(strings.Split(pth, "\\")...), j.Downloads.Artifact.Url); err != nil {
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
					if !checkSum(pth, classifier.Sha1) {
						log <- "downloader: downloading native '" + classifier.Path + "'"
						if err := getLib(path.Join(strings.Split(pth, "\\")...), classifier.Url); err != nil {
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

func getLib(out, url string) error {
	outDir := path.Dir(out)
	if err := os.MkdirAll(outDir, os.ModePerm); err != nil {
		return err
	}

	_, err := DownloadFile(url, out)
	return err
}

func checkSum(path, sha1check string) bool {
	data, err := os.ReadFile(path)
	if err != nil {
		return false
	}

	sha1Engine := sha1.New()
	sha1Engine.Write(data)
	sha1Hash := hex.EncodeToString(sha1Engine.Sum(nil))

	return sha1Hash == sha1check
}
