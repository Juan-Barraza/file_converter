package csv

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"github.com/DeijoseDevelop/file_converter/json"
	"github.com/DeijoseDevelop/file_converter/utils"
	"sync"
)

type ReadConvertFunc func(string) ([]map[string]any, error)

func ConvertToCsv(path, to string) error {
	file, fileErr := utils.OpenOrCreateFile("export.csv")
	if fileErr != nil {
		return fmt.Errorf(fileErr.Error())
	}
	defer file.Close()

	bufferedWriter := bufio.NewWriter(file)
	writer := csv.NewWriter(bufferedWriter)
	defer writer.Flush()

	readConvertOptions := map[string]ReadConvertFunc{
		"json": json.ReadJson,
		"csv":  json.ReadJson,
		"xml":  json.ReadJson,
		"yaml": json.ReadJson,
	}

	var maps []map[string]any

	if readConvertFunc, ok := readConvertOptions[to]; ok {
		data, err := readConvertFunc(path)
		if err != nil {
			return fmt.Errorf("error decoding file: %s", err)
		}
		maps = data
	}

	flatData := utils.FlattenSliceMap(maps)

	var headers []string
	for key := range flatData[0] {
		headers = append(headers, key)
	}
	writer.Write(headers)

	numWorkers := 4
	jobs := make(chan map[string]string, len(flatData))
	results := make(chan []string, len(flatData))
	var wg sync.WaitGroup

	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go utils.GoRoutineWorker(&wg, jobs, results, headers)
	}

	go func() {
		for _, record := range flatData {
			jobs <- record
		}
		close(jobs)
	}()

	go func() {
		wg.Wait()
		close(results)
	}()

	for record := range results {
		writer.Write(record)
	}

	bufferedWriter.Flush()

	return nil
}
