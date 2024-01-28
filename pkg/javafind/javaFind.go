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

	var jvms, dirNames []string

	if downloader.OperatingSystem == "linux" {
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

		return nil
	} else if downloader.OperatingSystem == "windows" {
		jvms = []string{"C:\\Program Files\\Java", "C:\\Program Files (x86)\\Java"}
		dirNames = []string{fmt.Sprintf("jre-1.%d", version), fmt.Sprintf("jre-%d", version)}
	} else if downloader.OperatingSystem == "osx" {
		jvms = []string{"/usr/local/opt/", "/Library/Java/JavaVirtualMachines"}
		dirNames = []string{fmt.Sprintf("openjdk@%d", version), fmt.Sprintf("openjdk-%d.jdk", version)}
	}

	for _, javaDir := range jvms {
		stat, err := os.Stat(javaDir)
		if !os.IsNotExist(err) && stat.IsDir() {
			for _, jvmDir := range dirNames {
				fpth := filepath.Join(javaDir, jvmDir, "bin", "java")
				stat, err = os.Stat(fpth)
				if !os.IsNotExist(err) && !stat.IsDir() {
					ret := filepath.Join(javaDir, jvmDir)
					return &ret
				}
			}
		}
	}

	return nil
}
