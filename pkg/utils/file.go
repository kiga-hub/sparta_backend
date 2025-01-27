package utils

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

// GetFileName -
func GetFileName(dir string) (string, error) {
	if !IsFileExist(dir) {
		return "", fmt.Errorf("dir is not exist")
	}

	var fileName string
	dirList, err := os.ReadDir(dir)
	if err != nil {
		return "", err
	}
	for _, fi := range dirList {
		if !fi.IsDir() {
			fileName = fi.Name()
			break
		}
	}
	if fileName == "" {
		return "", fmt.Errorf("fileName is nil")
	}
	return fileName, nil

}

// GetFolderName get folder name
func GetFolderName(dir string) (string, error) {
	if !IsFileExist(dir) {
		return "", fmt.Errorf("dir is not exist")
	}

	var dirName string
	dirList, err := os.ReadDir(dir)
	if err != nil {
		return "", err
	}
	for _, fi := range dirList {
		if fi.IsDir() {
			dirName = fi.Name()
			break
		}
	}
	if dirName == "" {
		return "", fmt.Errorf("dirName is nil")
	}
	return dirName, nil
}

// MakeDirIfNotExist -
func MakeDirIfNotExist(dir string) error {
	if _, err := os.Stat(dir); err != nil {
		if os.IsNotExist(err) {
			if err := os.MkdirAll(filepath.Dir(dir), 0755); err != nil {
				return err
			}
		}
	}
	return nil
}

// IsFileExist checks if the given file exists.
func IsFileExist(path string) bool {
	if _, err := os.Stat(path); err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}
	return true
}

// HasPrefix checks if the given string has the given prefix.
func HasPrefix(s, prefix string) bool {
	return len(s) >= len(prefix) && s[:len(prefix)] == prefix
}

// CopyDir copies a directory and its contents to a new location.
func CopyDir(src, dst string) error {
	// Create the destination directory if it doesn't exist
	if err := os.MkdirAll(dst, 0755); err != nil {
		return err
	}

	// Get a list of files and directories in the source directory
	entries, err := os.ReadDir(src)
	if err != nil {
		return err
	}

	// Copy each file and directory to the destination directory
	for _, entry := range entries {
		srcPath := filepath.Join(src, entry.Name())
		dstPath := filepath.Join(dst, entry.Name())

		if entry.IsDir() {
			// Recursively copy subdirectories
			if err := CopyDir(srcPath, dstPath); err != nil {
				return err
			}
		} else {
			// Copy files
			if err := CopyFile(srcPath, dstPath); err != nil {
				return err
			}
		}
	}

	return nil
}

// CopyFile copies a file to a new location.
func CopyFile(src, dst string) error {
	// Open the source file
	srcFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	// Create the destination file
	dstFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer dstFile.Close()

	// Copy the contents of the source file to the destination file
	if _, err := io.Copy(dstFile, srcFile); err != nil {
		return err
	}

	// Set the permissions of the destination file to match the source file
	srcInfo, err := os.Stat(src)
	if err != nil {
		return err
	}
	return os.Chmod(dst, srcInfo.Mode())
}

// ClearDir -
func ClearDir(dir string) error {
	if !IsFileExist(dir) {
		return nil
	}
	return os.RemoveAll(dir)
}

// ListFilesWithPrefix returns a list of files with the given prefix in the given directory.
func ListFilesWithPrefix(dirPath string, prefix string) ([]string, error) {
	var files []string

	err := filepath.Walk(dirPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && strings.HasPrefix(info.Name(), prefix) {
			files = append(files, path)
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	return files, nil
}
