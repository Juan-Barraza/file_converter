package xlsx

import (
	"fmt"

	"github.com/DeijoseDevelop/file_converter/converter"
	"github.com/DeijoseDevelop/file_converter/utils"
	"github.com/xuri/excelize/v2"
)

func ConvertToXlsx(path, to string) error {
	file, fileErr := utils.OpenOrCreateFile("export.xlsx")
	if fileErr != nil {
		return fmt.Errorf(fileErr.Error())
	}
	defer file.Close()
	
	xlsxCreate := excelize.NewFile()
	sheetName := "Sheet1"
	xlsxCreate.SetSheetName(xlsxCreate.GetSheetName(1), sheetName)

	var  maps []map[string]interface{}

	if readConverteFunc, ok := converter.GetReadConvertFunc(to); ok  {
		data, err := readConverteFunc(path)
		if err != nil {
			return fmt.Errorf("error decoding file %s", err)
		}
		maps = data
	}

	flatData := utils.FlattenSliceMap(maps)

	if len(flatData) > 0 {
		headers := make([]string, 0, len(flatData[0]))
		for key := range flatData[0] {
			headers = append(headers, key)
		}
		for i, header := range headers {
			cell, _ := excelize.CoordinatesToCellName(i+1, 1)
			xlsxCreate.SetCellValue(sheetName, cell, header)
		}
		for rowIdx, row := range flatData {
			for colIdx, header := range headers {
				cell, _ := excelize.CoordinatesToCellName(colIdx+1, rowIdx+2)
				xlsxCreate.SetCellValue(sheetName, cell, row[header])
			}
		}
	}
	if err := xlsxCreate.Write(file); err != nil {
		return fmt.Errorf("error writing to xlsx file: %s", err)
	}

	return nil
}