package javafind

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/zapomnij/firecraft/pkg/downloader"
)

func FindJava(version uint) *string {
	if version == 0 {
		version = 8
	}

	if downloader.OperatingSystem == "linux" {
		var jvms, dirnames []string

		_, err := os.Stat("/usr/bin/java")
		if !os.IsNotExist(err) {
			prelink, err := os.Readlink("/usr/bin/java")
			if err != nil || prelink == "/usr/bin/java" {
				return nil
			}

			if strings.HasSuffix(prelink, "/jre/bin/java") {
				jvms = append(dirnames, path.Dir(prelink)+"/../../../")
			} else if strings.HasSuffix(prelink, "/bin/java") {
				jvms = append(dirnames, path.Dir(prelink)+"/../../")
			}
		}

		if len(jvms) == 0 {
			jvms = []string{
				"/usr/lib/jvm",
				"/usr/lib64/jvm",
			}
		}

		dirnames = []string{
			fmt.Sprintf("java-%d-openjdk", version),
			fmt.Sprintf("java-1.%d.0-openjdk-1.%d.0", version, version),
			fmt.Sprintf("jre-%d", version),
			fmt.Sprintf("jre-%d-openjdk", version),
			fmt.Sprintf("jre-1.%d.0", version),
			fmt.Sprintf("jre-1.%d.0-openjdk", version),
		}

		for _, jvmsdir := range jvms {
			stat, err := os.Stat(jvmsdir)
			if !os.IsNotExist(err) && stat.IsDir() {
				for _, jvmdir := range dirnames {
					stat, err = os.Stat(jvmsdir + "/" + jvmdir)
					if !os.IsNotExist(err) && stat.IsDir() {
						javapaths := []string{
							path.Join(jvmsdir, jvmdir, "bin", "java"),
							path.Join(jvmsdir, jvmdir, "jre", "bin", "java"),
						}

						for _, jvbin := range javapaths {
							stat, err := os.Stat(jvbin)
							if !os.IsNotExist(err) && !stat.IsDir() {
								return &jvbin
							}
						}
					}
				}
			}
		}
	} else if downloader.OperatingSystem == "windows" {
		dirnames := []string{"C:\\Program Files\\Java", "C:\\Program Files (x86_64)\\Java"}
		for _, jvDir := range dirnames {
			stat, err := os.Stat(jvDir)
			if !os.IsNotExist(err) && stat.IsDir() {
				javaDirs := []string{filepath.Join(jvDir, fmt.Sprintf("jre-1.%d", version)), filepath.Join(jvDir, fmt.Sprintf("jre-%d", version))}
				for _, vmDir := range javaDirs {
					stat, err := os.Stat(vmDir)
					if !os.IsNotExist(err) && stat.IsDir() {
						ret := filepath.Join(vmDir, "bin", "java.exe")
						return &ret
					}
				}
			}
		}
	}

	return nil
}
