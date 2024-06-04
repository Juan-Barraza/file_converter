package utils

import (
	"fmt"
	"os"
)

func ReadFile(path string) (*os.File, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("error al abrir el archivo: %s", err)
	}
	return file, nil
}