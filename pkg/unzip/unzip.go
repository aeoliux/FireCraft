package unzip

import (
	"archive/zip"
	"io"
	"os"
	"path"
)

func Unzip(pth, dst string, exclude []string) error {
	archive, err := zip.OpenReader(pth)
	if err != nil {
		return err
	}
	defer archive.Close()

	for _, f := range archive.File {
		filedst := path.Join(dst, f.Name)
		if isArrayContaining(exclude, f.Name) {
			continue
		}

		if f.FileInfo().IsDir() {
			os.MkdirAll(filedst, os.ModePerm)
		} else {
			if err := os.MkdirAll(path.Dir(filedst), os.ModePerm); err != nil {
				return err
			}

			dstF, err := os.OpenFile(filedst, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
			if err != nil {
				return err
			}

			fileInArchive, err := f.Open()
			if err != nil {
				return err
			}

			if _, err := io.Copy(dstF, fileInArchive); err != nil {
				return err
			}

			dstF.Close()
			fileInArchive.Close()
		}
	}

	return nil
}
