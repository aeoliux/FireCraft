package runner

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/zapomnij/firecraft/pkg/downloader"
)

type Runner struct {
	Classpath string
	Username  string

	Version    downloader.VersionJSON
	JavaBinary string
	JavaArgs   string
	AssetIndex downloader.AssetIndex

	Xuid              *string
	Uuid              *string
	AccessToken       *string
	HaveBoughtTheGame bool
}

func NewRunner(username, javabinary, classpath string, javaargs string, verjson downloader.VersionJSON, assetIndex downloader.AssetIndex) *Runner {
	return &Runner{
		Username:   username,
		Version:    verjson,
		JavaBinary: javabinary,
		JavaArgs:   javaargs,
		Classpath:  classpath,
		AssetIndex: assetIndex,
	}
}

func (r *Runner) SetUpMicrosoft(uuid, accessToken string, haveBoughtTheGame bool) {
	r.Uuid = &uuid
	r.AccessToken = &accessToken
	r.HaveBoughtTheGame = haveBoughtTheGame
}

func (r Runner) parseJVMArg(arg string) string {
	arg = strings.ReplaceAll(arg, "${natives_directory}", downloader.NativesDir)
	arg = strings.ReplaceAll(arg, "${launcher_name}", "firecraft")
	arg = strings.ReplaceAll(arg, "${launcher_version}", "1")
	arg = strings.ReplaceAll(arg, "${classpath}", r.Classpath)

	return arg
}

func (r Runner) parseMCArg(arg string) string {
	// arg = strings.ReplaceAll(arg, "${auth_player_name}", r.Username)
	// arg = strings.ReplaceAll(arg, "${version_name}", r.Version.Id)
	// arg = strings.ReplaceAll(arg, "${game_directory}", downloader.MinecraftDir)
	// arg = strings.ReplaceAll(arg, "${assets_index_name}", downloader.AssetsDir)
	// arg = strings.ReplaceAll(arg, "${assets_root}", r.Version.AssetIndex.Id)
	// if r.Uuid != nil {
	// 	arg = strings.ReplaceAll(arg, "${auth_uuid}", r.Version.AssetIndex.Id)
	// }
	// if r.AccessToken != nil {
	// 	arg = strings.ReplaceAll(arg, "${auth_access_token}", *r.AccessToken)
	// 	arg = strings.ReplaceAll(arg, "${user_type}", "msa")
	// }
	// if r.Xuid != nil {
	// 	arg = strings.ReplaceAll(arg, "${auth_xuid}", *r.Xuid)
	// }
	// arg = strings.ReplaceAll(arg, "${version_type}", r.Version.Type)

	switch arg {
	case "${auth_player_name}":
		return r.Username
	case "${version_name}":
		return r.Version.Id
	case "${game_directory}":
		return downloader.MinecraftDir
	case "${assets_index_name}":
		return r.Version.Assets
	case "${assets_root}", "${game_assets}":
		if r.AssetIndex.Virtual {
			return filepath.Join(downloader.AssetsDir, "virtual", r.Version.Assets)
		} else if r.AssetIndex.MapToResources {
			return filepath.Join(downloader.MinecraftDir, "resources")
		}
		return downloader.AssetsDir
	case "${auth_uuid}":
		if r.Uuid == nil {
			return "null"
		}
		return *r.Uuid
	case "${auth_access_token}", "${auth_session}":
		if r.AccessToken == nil {
			return "null"
		}
		return *r.AccessToken
	case "${clientid}":
		return "null"
	case "${auth_xuid}":
		if r.Xuid == nil {
			return "null"
		}
		return *r.Xuid
	case "${user_type}":
		if r.AccessToken != nil {
			return "msa"
		} else {
			return "null"
		}
	case "${version_type}":
		return r.Version.Type
	}

	return arg
}

func (r Runner) Run() error {
	cmd := []string{r.JavaBinary}
	cmd = append(cmd, strings.Split(r.JavaArgs, " ")...)

	if r.Version.MinecraftArguments != nil {
		cmd = append(cmd, r.parseJVMArg("-Djava.library.path=${natives_directory}"), "-cp", r.Classpath, r.Version.MainClass)
		for _, word := range strings.Split(*r.Version.MinecraftArguments, " ") {
			cmd = append(cmd, r.parseMCArg(word))
		}
	} else {
		for _, jvmarg := range r.Version.Arguments.Jvm {
			switch jvmarg.(type) {
			case string:
				cmd = append(cmd, r.parseJVMArg(jvmarg.(string)))
			default:
				casted := jvmarg.(map[string]interface{})
				passed := false
				for _, rule := range casted["rules"].([]interface{}) {
					castedRule := rule.(map[string]interface{})
					if castedRule["action"].(string) == "disallow" {
						passed = true
					}

					if castedRule["name"] != nil {
						if castedRule["name"].(string) == downloader.OperatingSystem {
							passed = !passed
						}
					}
				}

				if passed {
					cmd = append(cmd, casted["value"].([]string)...)
				}
			}
		}

		cmd = append(cmd, r.Version.MainClass)

		for _, mcarg := range r.Version.Arguments.Game {
			switch mcarg.(type) {
			case string:
				cmd = append(cmd, r.parseMCArg(mcarg.(string)))
			}
		}
	}

	if !r.HaveBoughtTheGame && r.AccessToken != nil {
		cmd = append(cmd, "--demo")
	}

	run := exec.Command(cmd[0], cmd[1:]...)
	run.Stdout = os.Stdout
	run.Stderr = os.Stderr
	if err := run.Run(); err != nil {
		return err
	}

	return nil
}
