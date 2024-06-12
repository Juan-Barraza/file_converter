package yaml

import (
	"fmt"

	"github.com/DeijoseDevelop/file_converter/converter"
	"github.com/DeijoseDevelop/file_converter/utils"
	"gopkg.in/yaml.v2"
)

func ConvertToYaml(path, to string) error {
	file, fileErr := utils.OpenOrCreateFile("export.yaml")
	if fileErr != nil {
		return fmt.Errorf(fileErr.Error())
	}
	defer file.Close()

	encoder := yaml.NewEncoder(file)
	defer encoder.Close()

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
		return fmt.Errorf("error writing YAML data: %s", err)
	}

	return nil
}
