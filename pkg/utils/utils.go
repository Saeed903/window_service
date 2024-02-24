package utils

import (
	"os"
	"path/filepath"
)

// ReadFiles read from directory and get path all file
func ReadFiles(folderPath string) ([]string, error) {

	var fileList []string

	// Read path folder
	files, err := os.ReadDir(folderPath)
	if err != nil {
		return nil, err
	}

	// explore all entire files and subfolders
	for _, file := range files {
		if file.IsDir() {
			subFolderpath := filepath.Join(folderPath, file.Name())
			subFileList, err := ReadFiles(subFolderpath)
			if err != nil {
				return nil, err
			}
			fileList = append(fileList, subFileList...)
		}else {
			// check file is json
			filePath := filepath.Join(folderPath, file.Name())
			fileList = append(fileList, filePath)
		}
	} 
	return fileList, nil
}

