// Author: Rui
// Date: 2023/01/08 16:20
// Description: calculate folder total size

package util

import (
	"os"
	"path/filepath"
)

func GetDirSize(dirPath string) (int64, error) {
	var size int64 = 0

	err := filepath.Walk(dirPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		size += info.Size()

		return nil
	})

	if err != nil {
		return 0, err
	}

	return size, nil
}
