package window

import (
	"encoding/json"
	"os"
	"path"

	"github.com/zapomnij/firecraft/pkg/downloader"
)

type LProfile struct {
	Name          string `json:"name"`
	Type          string `json:"type"`
	LastVersionId string `json:"lastVersionId"`
	GameDir       string `json:"gameDir"`
	JavaBin       string `json:"javaBin"`
	JavaArgs      string `json:"javaArgs"`
}

type LauncherProfiles struct {
	Profiles               map[string]LProfile `json:"profiles"`
	AuthenticationDatabase struct {
		RefreshToken string
		Email        string
		Username     string
	}
}

func (l LauncherProfiles) Save() error {
	pth := path.Join(downloader.MinecraftDir, "launcher_profiles.json")

	bt, err := json.Marshal(l)
	if err != nil {
		return err
	}

	if err := os.WriteFile(pth, bt, os.ModePerm); err != nil {
		return err
	}

	return nil
}

func loadProfiles() (*LauncherProfiles, error) {
	pth := path.Join(downloader.MinecraftDir, "launcher_profiles.json")
	rd, err := os.ReadFile(pth)
	if err != nil {
		if os.IsNotExist(err) {
			l := LauncherProfiles{
				Profiles: make(map[string]LProfile),
			}

			return &l, nil
		}

		return nil, err
	}

	var l LauncherProfiles
	if err := json.Unmarshal(rd, &l); err != nil {
		return nil, err
	}

	return &l, nil
}

func (l *LauncherProfiles) deleteProfile(name string) {
	delete(l.Profiles, name)
}
