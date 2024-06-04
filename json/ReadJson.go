package json

import (
	"encoding/json"
	"fmt"
	"github/file_converter/utils"
)

func ReadJson(path string) ([]map[string]any, error) {
	file, err := utils.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("error opening the file: %s", err)
	}
	defer file.Close()

	var data any
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&data)

	if err != nil {
		return nil, fmt.Errorf("error decoding the JSON: %s", err)
	}

	maps := make([]map[string]any, 0)

	switch currentData := data.(type) {
	case map[string]any:
		fmt.Println("The JSON is an object:")
		maps = append(maps, currentData)
	case []any:
		fmt.Println("The JSON is a list:")
		for _, value := range currentData {
			if mapValue, ok := value.(map[string]any); ok {
				maps = append(maps, mapValue)
			}
		}
	default:
		return nil, fmt.Errorf("the json has an unknow format")
	}

	return maps, nil
}
