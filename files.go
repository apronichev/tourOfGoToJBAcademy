package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
)

// DeleteContents removes all files and subdirectories in a directory
func DeleteContents(dir string) error {
	// List all the contents of the directory
	files, err := os.ReadDir(dir)
	if err != nil {
		return err
	}

	// Iterate over each entry and remove it
	for _, file := range files {
		filePath := filepath.Join(dir, file.Name())
		if file.IsDir() {
			// If the entry is a directory, delete its contents first
			err := DeleteContents(filePath)
			if err != nil {
				return err
			}
			// Remove the directory
			os.Remove(filePath)
		} else {
			// Remove the file
			os.Remove(filePath)
		}
	}
	return nil
}

func FindFile(root, filename string) (string, error) {
	var foundPath string
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		if filepath.Base(path) == filename {
			foundPath = path
			fmt.Println("File found: ", path)
			return io.EOF // Using io.EOF to signal that the file was found
		}
		return nil
	})

	if err != nil {
		if err == io.EOF {
			return foundPath, nil // File found, return the path
		}
		return "", err // Some other error occurred
	}
	return "", fmt.Errorf("file not found") // File not found
}

func FindArticleFiles(directory string) []string {
	// List of filenames to search for, in the required order
	filesToFind := map[string]string{
		"basics.article":      "",
		"flowcontrol.article": "",
		"moretypes.article":   "",
		"methods.article":     "",
		"generics.article":    "",
		"concurrency.article": "",
	}

	// Walk through the directory to find the files
	filepath.Walk(directory, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		filename := filepath.Base(path)
		if _, exists := filesToFind[filename]; exists {
			filesToFind[filename] = path
		}
		return nil
	})

	// Create the foundFiles slice in the required order
	foundFiles := []string{
		filesToFind["basics.article"],
		filesToFind["flowcontrol.article"],
		filesToFind["moretypes.article"],
		filesToFind["methods.article"],
		filesToFind["generics.article"],
		filesToFind["concurrency.article"],
	}

	return foundFiles
}
