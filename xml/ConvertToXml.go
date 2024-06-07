package xml

import (
	"encoding/xml"
	"fmt"
	"sync"

	myJson "github.com/DeijoseDevelop/file_converter/json"
	"github.com/DeijoseDevelop/file_converter/utils"
	csv "github.com/DeijoseDevelop/file_converter/csv"

	
)

type MapEntry struct {
	XMLName xml.Name
	Value   string `xml:",chardata"`
}

type ReadConvertFunc func(string) ([]map[string]any, error)

func ConvertToXml(path, to string) error {
	file, fileErr := utils.OpenOrCreateFile("export.xml")
	if fileErr != nil {
		return fmt.Errorf(fileErr.Error())
	}
	defer file.Close()

	readConvertOptions := map[string]ReadConvertFunc{
		"json": myJson.ReadJson,
		"csv":  csv.ReadCSV, 
		"xml":  myJson.ReadJson, 
		"yaml": myJson.ReadJson, 
	}

	var maps []map[string]any

	if readConvertFunc, ok := readConvertOptions[to]; ok {
		data, err := readConvertFunc(path)
		if err != nil {
			return fmt.Errorf("error decoding file: %s", err)
		}
		maps = data
	} else {
		return fmt.Errorf("unsupported format: %s", to)
	}

	flatData := utils.FlattenSliceMap(maps)

	numWorkers := 4
	jobs := make(chan map[string]string, len(flatData))
	results := make(chan []byte, len(flatData))
	var wg sync.WaitGroup

	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for record := range jobs {
				entries := make([]MapEntry, 0, len(record))
				for k, v := range record {
					entries = append(entries, MapEntry{XMLName: xml.Name{Local: k}, Value: fmt.Sprintf("%v", v)})
				}
				itemXml, err := xml.MarshalIndent(entries, "", "   ")
				if err != nil {
					fmt.Printf("error converting item to XML: %s\n", err)
					continue
				}
				results <- itemXml
			}
		}()
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

	for itemXml := range results {
		if _, err := file.Write(itemXml); err != nil {
			return fmt.Errorf("error writing XML file: %s", err)
		}
		if _, err := file.Write([]byte("\n")); err != nil {
			return fmt.Errorf("error writing new line to XML file: %s", err)
		}
	}

	return nil
}
