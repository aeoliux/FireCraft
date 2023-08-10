package downloader

import (
	"encoding/json"
	"fmt"
	"os"
	"path"
	"strings"
)

type AssetIndex struct {
	Objects map[string]struct {
		Hash string
		Size uint
	}

	Virtual        bool `json:"virtual,omitempty"`
	MapToResources bool `json:"map_to_resources,omitempty"`
}

var Ai AssetIndex

func (v VersionJSON) GetAssets(log chan string) error {
	assetIndexPath := path.Join(AssetsDir, "indexes", v.Assets+".json")
	if err := os.MkdirAll(path.Dir(assetIndexPath), os.ModePerm); err != nil {
		log <- "error: " + err.Error()
		return err
	}

	if v.AssetIndex == nil {
		log <- "downloader: assets finished"
		return nil
	}

	log <- "downloader: downloading asset index"
	img, err := DownloadFile(v.AssetIndex.Url, assetIndexPath)
	if err != nil {
		log <- "error: " + err.Error()
		return err
	}

	if err := json.Unmarshal(img, &Ai); err != nil {
		log <- "error: " + err.Error()
		return err
	}

	var objDir string
	if Ai.Virtual {
		objDir = path.Join(AssetsDir, "virtual", v.Assets)
	} else if Ai.MapToResources {
		objDir = path.Join(MinecraftDir, "resources")
	} else {
		objDir = path.Join(AssetsDir, "objects")
	}
	for k, v := range Ai.Objects {
		var outDir string
		if Ai.Virtual || Ai.MapToResources {
			split := []string{objDir}
			split = append(split, strings.Split(path.Dir(k), "/")...)
			outDir = path.Join(split...)
		} else {
			outDir = path.Join(objDir, v.Hash[:2])
		}
		if err := os.MkdirAll(outDir, os.ModePerm); err != nil {
			log <- "error: " + err.Error()
			return err
		}

		var out string
		if Ai.Virtual || Ai.MapToResources {
			out = path.Join(outDir, path.Base(k))
		} else {
			out = path.Join(outDir, v.Hash)
		}
		if !checkSum(out, v.Hash) {
			log <- "downloader: downloading " + k
			url := fmt.Sprintf("https://resources.download.minecraft.net/%s/%s", v.Hash[:2], v.Hash)
			if _, err = DownloadFile(url, out); err != nil {
				log <- "error: " + err.Error()
				return err
			}
		}
	}

	log <- "downloader: assets finished"
	return nil
}
