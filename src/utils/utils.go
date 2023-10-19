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

func CreateNewFile(dirPath string, name string) {
	file, _ := os.OpenFile(path.Join(dirPath, name), os.O_CREATE, os.ModePerm)
	file.Close()
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
