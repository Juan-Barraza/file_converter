package json

import (
	"encoding/json"
	"fmt"
	"github.com/DeijoseDevelop/file_converter/converter"
	"github.com/DeijoseDevelop/file_converter/utils"
)

func init() {
	converter.RegisterReadConvertFunc("json", ReadJson)
}

func ReadJson(path string) ([]map[string]interface{}, error) {
	file, err := utils.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("error opening the file: %s", err)
	}
	defer file.Close()

	var data interface{}
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&data)

	if err != nil {
		return nil, fmt.Errorf("error decoding the JSON: %s", err)
	}

	maps := make([]map[string]interface{}, 0)

	switch currentData := data.(type) {
	case map[string]interface{}:
		fmt.Println("The JSON is an object:")
		maps = append(maps, currentData)
	case []interface{}:
		fmt.Println("The JSON is a list:")
		for _, value := range currentData {
			if mapValue, ok := value.(map[string]interface{}); ok {
				maps = append(maps, mapValue)
			}
		}
	default:
		return nil, fmt.Errorf("the JSON has an unknown format")
	}


	fmt.Println("Datos decodificados del JSON:", maps)

	return maps, nil
}
