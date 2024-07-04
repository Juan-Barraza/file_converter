package xlsx

import (
	"fmt"

	"github.com/DeijoseDevelop/file_converter/utils"
	"github.com/xuri/excelize/v2"
)

func ReadXlsx(path string) ([]map[string]interface{}, error){
	file, err := utils.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("error opening file: %s", err)
	}
	defer file.Close()
	
	xlsx, err := excelize.OpenReader(file)
	if err != nil {
		return nil, fmt.Errorf("error opening xlsx file: %s", err)
	}

	sheetCount := xlsx.SheetCount
	if sheetCount < 1 {
		return nil, fmt.Errorf("no sheets found in the Excel file")
	}

	sheetNames := xlsx.GetSheetList()
	if len(sheetNames) == 0 {
		return nil, fmt.Errorf("no sheets found in the Excel file")
	}
	sheetName := sheetNames[0]


	rows, err := xlsx.GetRows(sheetName)
	if err != nil {
		return nil, fmt.Errorf("error reading rows from sheet: %s", err)
	}

	if len(rows)  < 1 {
		return nil, fmt.Errorf("no rows found in the sheet")
	}

	headers := rows[0]
	var result []map[string]interface{}

	for _, row := range rows[1:] {
		rowMap := make(map[string]interface{})
		for i, cell := range row {
			if i < len(headers){
				rowMap[headers[i]] = cell
			}
		}
		result = append(result,  rowMap)
	}


	return result, nil
}