package utils

import (
	"crypto/rand"
	"encoding/hex"
	"io"
	"os"
	"path"
	"strings"
)

func GetMapKeys[K comparable, V any](value map[K]V) []K {
	result := []K{}
	for key := range value {
		result = append(result, key)
	}
	return result
}

func GenerateSecureToken(length int) string {
	token := make([]byte, length)
	if _, err := rand.Read(token); err != nil {
		return ""
	}
	return hex.EncodeToString(token)
}

func FormatTokemizedEndpoint(endpoint string, token string) string {
	return strings.NewReplacer("{token}", token).Replace(endpoint)
}

func CreateNewDirectory(dirPath string, name string) {
	os.MkdirAll(path.Join(dirPath, name), os.ModePerm)
}

func CreateNewFile(dirPath string, name string) error {
	file, err := os.OpenFile(path.Join(dirPath, name), os.O_CREATE, os.ModePerm)
	if err != nil {
		return err
	}
	file.Close()
	return nil
}

func RemoveFile(dirPath string) error {
	return os.RemoveAll(dirPath)
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

func Rename(srcPath string, oldName string, newName string) error {
	return os.Rename(path.Join(srcPath, oldName), path.Join(srcPath, newName))
}

func HandleFilesystemPath(value string) string {
	value = strings.Trim(value, "\"")
	value = path.Clean(value)
	return value
}

func IfNotNil[T any](object *T, handler func(object *T) error) error {
	if object != nil {
		return handler(object)
	}
	return nil
}
