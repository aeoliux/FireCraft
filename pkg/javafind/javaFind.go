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
	// Default version is 8
	if version == 0 {
		version = 8
	}

	if downloader.OperatingSystem == "linux" {
		var jvms, dirNames []string

		// Find system JVMs directory
		_, err := os.Stat("/usr/bin/java")
		if !os.IsNotExist(err) {
			prelink, err := os.Readlink("/usr/bin/java")
			if err == nil || prelink != "/usr/bin/java" {
				if strings.HasSuffix(prelink, "/jre/bin/java") {
					// If it is JRE, not JDK
					jvms = append(jvms, path.Dir(prelink)+"/../../../")
				} else if strings.HasSuffix(prelink, "/bin/java") {
					// When JDK
					jvms = append(jvms, path.Dir(prelink)+"/../../")
				}
			}
		}

		// Most common JVM directories
		jvms = append(jvms, "/usr/lib/jvm", "/usr/lib64/jvm")

		// Common directory names for JVMs
		dirNames = []string{
			// JDK: Arch Linux and probably Ubuntu
			fmt.Sprintf("java-%d-openjdk", version),
			// JDK: OpenSUSE
			fmt.Sprintf("java-1.%d.0-openjdk-1.%d.0", version, version),
			fmt.Sprintf("jre-%d", version),
			fmt.Sprintf("jre-%d-openjdk", version),
			fmt.Sprintf("jre-1.%d.0", version),
			fmt.Sprintf("jre-1.%d.0-openjdk", version),
		}

		// Check for defined JVM paths
		for _, jvmsDir := range jvms {
			stat, err := os.Stat(jvmsDir)
			if !os.IsNotExist(err) && stat.IsDir() {
				// Check for possible JVM directories
				for _, jvmDir := range dirNames {
					stat, err = os.Stat(jvmsDir + "/" + jvmDir)
					if !os.IsNotExist(err) && stat.IsDir() {
						// Possible JVM binary paths
						javaPaths := []string{
							path.Join(jvmsDir, jvmDir, "bin", "java"),
							path.Join(jvmsDir, jvmDir, "jre", "bin", "java"),
						}

						for _, jvBin := range javaPaths {
							stat, err := os.Stat(jvBin)
							if !os.IsNotExist(err) && !stat.IsDir() {
								// If it exists, just get the path of JVM, not path to binary
								jvDir := path.Dir(path.Dir(jvBin))
								return &jvDir
							}
						}
					}
				}
			}
		}
	} else if downloader.OperatingSystem == "windows" {
		dirNames := []string{"C:\\Program Files\\Java", "C:\\Program Files (x86)\\Java"}
		for _, jvDir := range dirNames {
			stat, err := os.Stat(jvDir)
			if !os.IsNotExist(err) && stat.IsDir() {
				javaDirs := []string{filepath.Join(jvDir, fmt.Sprintf("jre-1.%d", version)), filepath.Join(jvDir, fmt.Sprintf("jre-%d", version))}
				for _, vmDir := range javaDirs {
					stat, err := os.Stat(vmDir)
					if !os.IsNotExist(err) && stat.IsDir() {
						ret := vmDir
						return &ret
					}
				}
			}
		}
	} else if downloader.OperatingSystem == "osx" {

	}

	return nil
}
