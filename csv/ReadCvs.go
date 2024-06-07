package csv

import (
	"encoding/csv"
	"fmt"

	"github.com/DeijoseDevelop/file_converter/utils"
)

type ReadConvert func(string) ([]map[string]any, error)

func ReadCSV(path string) ([]map[string]any, error) {
	file, err := utils.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("error opening the file: %s", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	headers, err := reader.Read()
	if err != nil {
		return nil, fmt.Errorf("error reading the csv headers: %s", err)
	}

	var records []map[string]any

	for {
		record, err := reader.Read()
		if err != nil {
			break
		}

		row := make(map[string]any)
		for i, header := range headers {
			row[header] = record[i]
		}

		records = append(records, row)
	}


	return records, nil
}
