package utils

import (
	"io"
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

func CopyFile(srcPath string, dstPath string) error {
	in, err := os.Open(srcPath)
	if err != nil {
		return err
	}
	defer in.Close()

	out, err := os.Create(dstPath)
	if err != nil {
		return err
	}

	defer func() {
		cerr := out.Close()
		if err == nil {
			err = cerr
		}
	}()

	if _, err = io.Copy(out, in); err != nil {
		return err
	}
	err = out.Sync()
	return nil
}
