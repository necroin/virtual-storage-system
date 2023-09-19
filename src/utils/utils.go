package utils

import (
	"os"
	"path"
)

func GetMapKeys[K comparable, V any](value map[K]V) []K {
	result := []K{}
	for key := range value {
		result = append(result, key)
	}
	return result
}

func CreateNewDirectory(dirPath string) {
	os.MkdirAll(path.Join(dirPath, "Новая папка"), os.ModePerm)
}

func RemoveDirectory(dirPath string) {
	os.RemoveAll(dirPath)
}
