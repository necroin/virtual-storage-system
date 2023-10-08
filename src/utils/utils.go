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

func CreateNewDirectory(dirPath string, name string) {
	os.MkdirAll(path.Join(dirPath, name), os.ModePerm)
}

func CreateNewFile(dirPath string, name string) {
	file, _ := os.OpenFile(path.Join(dirPath, name), os.O_CREATE, os.ModePerm)
	file.Close()
}

func RemoveFile(dirPath string) {
	os.RemoveAll(dirPath)
}
