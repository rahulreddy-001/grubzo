package utils

import (
	"encoding/json"
	"os"
)

func LoadJSONFromFile[T any](filePath string) (*T, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var value T

	if err := json.NewDecoder(file).Decode(&value); err != nil {
		return nil, err
	}

	return &value, nil
}
