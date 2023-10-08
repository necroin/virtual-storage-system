package buffer

import "fmt"

var (
	filePath string
)

func SetFile(path string) {
	filePath = path
}

func GetFile() (string, error) {
	if filePath == "" {
		return filePath, fmt.Errorf("[Buffer] [Error] buffered file path is empty")
	}
	return filePath, nil
}
