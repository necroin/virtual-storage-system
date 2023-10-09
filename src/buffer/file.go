package buffer

import "fmt"

var (
	bufFilePath string
	bufFileType string
)

func SetFile(filePath string, fileType string) {
	bufFilePath = filePath
	bufFileType = fileType
}

func GetFile() (string, string, error) {
	if bufFilePath == "" {
		return bufFilePath, bufFileType, fmt.Errorf("[Buffer] [Error] buffered file path is empty")
	}
	return bufFilePath, bufFileType, nil
}
