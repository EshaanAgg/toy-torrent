package utils

import (
	"fmt"
	"os"
)

func MakeFileWithData(filePath string, data []byte) error {
	// Create the file with all the parent directories
	err := os.MkdirAll(filePath, os.ModePerm)
	if err != nil {
		return fmt.Errorf("error creating directories for file: %w", err)
	}
	file, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("error creating file: %w", err)
	}
	defer file.Close()

	// Write the data to the file
	_, err = file.Write(data)
	if err != nil {
		return fmt.Errorf("error writing data to file: %w", err)
	}
	// Flush the data to disk
	err = file.Sync()
	if err != nil {
		return fmt.Errorf("error flushing data to disk: %w", err)
	}

	return nil
}
