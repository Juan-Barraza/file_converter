package utils

import "fmt"

func FlattenSliceMap(data []map[string]interface{}) []map[string]string {
	flatData := make([]map[string]string, len(data))
	for i, item := range data {
		flatData[i] = make(map[string]string)
		flattenMap("", item, flatData[i])
	}
	return flatData
}

func flattenMap(prefix string, input map[string]interface{}, output map[string]string) {
	for key, value := range input {
		fullKey := key
		if prefix != "" {
			fullKey = key
		}
		switch value := value.(type) {
		case map[string]any:
			flattenMap(fullKey, value, output)
		default:
			output[fullKey] = fmt.Sprint(value)
		}
	}
}