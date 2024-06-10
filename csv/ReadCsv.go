package csv

import (
	"encoding/csv"
	"fmt"
	"os"
	"github.com/DeijoseDevelop/file_converter/converter"
)

func init() {
	converter.RegisterReadConvertFunc("csv", ReadCSV)
}

func ReadCSV(path string) ([]map[string]interface{}, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("error opening the file: %s", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	headers, err := reader.Read()
	if err != nil {
		return nil, fmt.Errorf("error reading the csv headers: %s", err)
	}

	var records []map[string]interface{}
	for {
		record, err := reader.Read()
		if err != nil {
			break
		}

		row := make(map[string]interface{})
		for i, header := range headers {
			row[header] = record[i]
		}

		records = append(records, row)
	}

	return records, nil
}
