// Author: Rui
// Date: 2023/01/05 16:22
// Description: compress file toolkit

package util

import (
	"archive/zip"
	"io"
	"os"
	"path/filepath"
)

type CompressUtil struct {
}

func (p *CompressUtil) Compress(folderPath string, zipFileName string) error {
	zipFile, err := os.Create(zipFileName)
	if err != nil {
		return err
	}
	defer zipFile.Close()

	zipWriter := zip.NewWriter(zipFile)
	defer zipWriter.Close()

	err = filepath.Walk(folderPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.Mode().IsRegular() {
			return nil
		}

		file, err := os.Open(path)
		if err != nil {
			return err
		}
		defer file.Close()

		zipFile, err := zipWriter.Create(path[len(folderPath):])
		if err != nil {
			return err
		}
		_, err = io.Copy(zipFile, file)
		if err != nil {
			return err
		}

		return nil
	})

	return err
}
