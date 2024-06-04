package utils

import (
	"fmt"
	"os"
)

func OpenOrCreateFile(path string) (*os.File, error) {
	_, err := os.Stat(path)
	var file *os.File

	if err == nil {
		file, err = os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			return nil, fmt.Errorf("error opening existing file: %v", err)
		}
	} else if os.IsNotExist(err) {
		file, err = os.Create(path)
		if err != nil {
			return nil, fmt.Errorf("error when creating the file: %v", err)
		}
	} else {
		return nil, fmt.Errorf("error checking file existence: %v", err)
	}

	return file, nil
}
