package xml

import (
	"encoding/xml"
	"fmt"
	"sync"
	"github.com/DeijoseDevelop/file_converter/utils"
	"github.com/DeijoseDevelop/file_converter/converter"
)

type XmlMapEntry struct {
	XMLName xml.Name
	Value   string `xml:",chardata"`
}

type XmlRecord struct {
	XMLName xml.Name `xml:"record"`
	Entries []XmlMapEntry `xml:",any"`
}

type XmlRoot struct {
	XMLName xml.Name `xml:"records"`
	Records []XmlRecord `xml:"record"`
}

func ConvertToXml(path, to string) error {
	file, fileErr := utils.OpenOrCreateFile("export.xml")
	if fileErr != nil {
		return fmt.Errorf(fileErr.Error())
	}
	defer file.Close()

	readConvertFunc, ok := converter.GetReadConvertFunc(to)
	if !ok {
		return fmt.Errorf("unsupported conversion type: %s", to)
	}

	data, err := readConvertFunc(path)
	if err != nil {
		return fmt.Errorf("error decoding file: %s", err)
	}

	flatData := utils.FlattenSliceMap(data)

	numWorkers := 4
	jobs := make(chan map[string]string, len(flatData))
	results := make(chan XmlRecord, len(flatData))
	var wg sync.WaitGroup

	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for record := range jobs {
				entries := make([]XmlMapEntry, 0, len(record))
				for k, v := range record {
					entries = append(entries, XmlMapEntry{XMLName: xml.Name{Local: k}, Value: fmt.Sprintf("%v", v)})
				}
				results <- XmlRecord{Entries: entries}
			}
		}()
	}

	go func() {
		for _, record := range flatData {
			job := make(map[string]string)
			for k, v := range record {
				job[k] = fmt.Sprint(v)
			}
			jobs <- job
		}
		close(jobs)
	}()

	go func() {
		wg.Wait()
		close(results)
	}()

	var root XmlRoot
	for record := range results {
		root.Records = append(root.Records, record)
	}

	rootXml, err := xml.MarshalIndent(root, "", "   ")
	if err != nil {
		return fmt.Errorf("error marshaling XML: %s", err)
	}

	if _, err := file.Write(rootXml); err != nil {
		return fmt.Errorf("error writing XML file: %s", err)
	}

	return nil
}
