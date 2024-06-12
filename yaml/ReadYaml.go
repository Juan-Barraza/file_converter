package yaml

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"github.com/DeijoseDevelop/file_converter/utils"
)

func ReadYaml(path string) ([]map[string]any, error) {
	file, err := utils.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("error opening the file: %s", err)
	}
	defer file.Close()

	var data any
	decoder := yaml.NewDecoder(file)
	err = decoder.Decode(&data)

	if err != nil {
		return nil, fmt.Errorf("error decoding the YAML: %s", err)
	}

	maps := make([]map[string]any, 0)

	switch currentData := data.(type) {
	case map[any]any:
		fmt.Println("The YAML is an object:")
		maps = append(maps, convertMap(currentData))
	case []any:
		fmt.Println("The YAML is a list:")
		for _, value := range currentData {
			if mapValue, ok := value.(map[any]any); ok {
				maps = append(maps, convertMap(mapValue))
			}
		}
	default:
		return nil, fmt.Errorf("the YAML has an unknown format")
	}

	return maps, nil
}

func convertMap(data map[any]any) map[string]any {
	result := make(map[string]any)
	for key, value := range data {
		strKey := fmt.Sprintf("%v", key)
		result[strKey] = value
	}
	return result
}
