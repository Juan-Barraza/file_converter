package csv

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"os"
	"sync"

	"github.com/DeijoseDevelop/file_converter/converter"
	"github.com/DeijoseDevelop/file_converter/utils"
)

func ConvertToCsv(path, to string) error {
	file, fileErr := os.Create("export.csv")
	if fileErr != nil {
		return fmt.Errorf(fileErr.Error())
	}
	defer file.Close()

	bufferedWriter := bufio.NewWriter(file)
	writer := csv.NewWriter(bufferedWriter)
	defer writer.Flush()

	var maps []map[string]interface{}

	if readConvertFunc, ok := converter.GetReadConvertFunc(to); ok {
		data, err := readConvertFunc(path)
		if err != nil {
			return fmt.Errorf("error decoding file: %s", err)
		}
		maps = data
	}

	flatData := utils.FlattenSliceMap(maps)

	fmt.Println("Datos decodificados del JSON:", maps)


	if len(flatData) == 0 {
		return fmt.Errorf("no data to convert to CSV")
	}

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
