package json

import (
	jsonPackage "encoding/json"
	"fmt"

	"github.com/DeijoseDevelop/file_converter/converter"
	"github.com/DeijoseDevelop/file_converter/utils"
)


func ConvertToJson(path, to string) error {
	file, fileErr := utils.OpenOrCreateFile("export.json")
	if fileErr != nil {
		return fmt.Errorf(fileErr.Error())
	}
	defer file.Close()

	encoder := jsonPackage.NewEncoder(file)
	encoder.SetIndent("", "   ")

	var maps []map[string]any

	if readConvertFunc, ok := converter.GetReadConvertFunc(to); ok {
		data, err := readConvertFunc(path)
		if err != nil {
			return fmt.Errorf("error decoding file: %s", err)
		}
		maps = data
	}

	flatData := utils.FlattenSliceMap(maps)

	err := encoder.Encode(flatData)
	if err != nil {
		return fmt.Errorf("error encoding data to JSON: %s", err)
	}

	return nil
}