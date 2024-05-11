package utils

import (
	"archive/tar"
	"compress/gzip"
	"fmt"
	"io"
	"os"
	"path"
	"path/filepath"
)

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

func Compress(archivePath string, buffer io.Writer) error {
	gzipWriter := gzip.NewWriter(buffer)
	tarWriter := tar.NewWriter(gzipWriter)

	fileInfo, err := os.Stat(archivePath)
	if err != nil {
		return fmt.Errorf("[Compress] failed get file info: %s", err)
	}

	mode := fileInfo.Mode()
	if mode.IsRegular() {
		header, err := tar.FileInfoHeader(fileInfo, archivePath)
		if err != nil {
			return fmt.Errorf("[Compress] failed get header by file info: %s", err)
		}

		if err := tarWriter.WriteHeader(header); err != nil {
			return fmt.Errorf("[Compress] failed write file header: %s", err)
		}

		data, err := os.Open(archivePath)
		if err != nil {
			return fmt.Errorf("[Compress] failed open file: %s", err)
		}
		defer data.Close()

		if _, err := io.Copy(tarWriter, data); err != nil {
			return fmt.Errorf("[Compress] failed copy file data: %s", err)
		}
	} else if mode.IsDir() {
		filepath.Walk(archivePath, func(file string, fileInfo os.FileInfo, err error) error {
			header, err := tar.FileInfoHeader(fileInfo, file)
			if err != nil {
				return fmt.Errorf("[Compress] failed get header by file info: %s", err)
			}

			relativePath, err := filepath.Rel(filepath.Dir(archivePath), file)
			if err != nil {
				return fmt.Errorf("[FindFilesInDir] failed get relative path: %s", err)
			}
			header.Name = filepath.ToSlash(relativePath)

			if err := tarWriter.WriteHeader(header); err != nil {
				return fmt.Errorf("[Compress] failed write file header: %s", err)
			}

			if !fileInfo.IsDir() {
				data, err := os.Open(file)
				if err != nil {
					return fmt.Errorf("[Compress] failed open file: %s", err)
				}
				defer data.Close()

				if _, err := io.Copy(tarWriter, data); err != nil {
					return fmt.Errorf("[Compress] failed copy file data: %s", err)
				}
			}
			return nil
		})
	} else {
		return fmt.Errorf("[Compress] file type not supported")
	}

	if err := tarWriter.Close(); err != nil {
		return fmt.Errorf("[Compress] failed close tar writer: %s", err)
	}

	if err := gzipWriter.Close(); err != nil {
		return fmt.Errorf("[Compress] failed close gzip writer: %s", err)
	}

	return nil
}

func Decompress(src io.Reader, dst string) error {
	if err := os.MkdirAll(dst, os.ModePerm); err != nil {
		return fmt.Errorf("[Decompress] failed create output dir: %s", err)
	}

	gzipReader, err := gzip.NewReader(src)
	if err != nil {
		return fmt.Errorf("[Decompress] failed get gzip reader: %s", err)
	}
	defer gzipReader.Close()

	tarReader := tar.NewReader(gzipReader)

	for {
		header, err := tarReader.Next()
		if err == io.EOF {
			break // End of archive
		}
		if err != nil {
			return err
		}
		target := filepath.Join(dst, header.Name)

		switch header.Typeflag {
		case tar.TypeDir:
			if _, err := os.Stat(target); err != nil {
				if err := os.MkdirAll(target, 0755); err != nil {
					return fmt.Errorf("[Decompress] failed create directory: %s", err)
				}
			}
		case tar.TypeReg:
			fileToWrite, err := os.OpenFile(target, os.O_CREATE|os.O_RDWR|os.O_TRUNC, os.FileMode(header.Mode))
			if err != nil {
				return fmt.Errorf("[Decompress] failed create file: %s", err)
			}
			defer fileToWrite.Close()

			if _, err := io.Copy(fileToWrite, tarReader); err != nil {
				return fmt.Errorf("[Decompress] failed copy file data: %s", err)
			}
		}
	}

	return nil
}
