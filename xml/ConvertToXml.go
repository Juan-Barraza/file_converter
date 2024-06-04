package xml

import (
	"encoding/xml"
	"fmt"
	myJson "github/file_converter/json"
	"github/file_converter/utils"
	"os"
)

type MapEntry struct {
	XMLName xml.Name
	Value   string `xml:",chardata"`
}

func ConvertToXml(path string) error {
	file, fileErr := utils.OpenOrCreateFile("export.xml")
	if fileErr != nil {
		return fmt.Errorf(fileErr.Error())
	}
	defer file.Close()

	jsonFile, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("error opening the file: %s", err)
	}
	defer jsonFile.Close()

	jsonData, err := myJson.ReadJson(path)
	if err != nil {
		return fmt.Errorf("error decoding file: %s", err)
	}

	for _, item := range utils.FlattenSliceMap(jsonData) {
		entries := make([]MapEntry, 0, len(item))
		for k, v := range item {
			entries = append(entries, MapEntry{XMLName: xml.Name{Local: k}, Value: fmt.Sprintf("%v", v)})
		}
		itemXml, err := xml.MarshalIndent(entries, "", "   ")
		if err != nil {
			return fmt.Errorf("error converting item to XML: %s", err)
		}
		if _, err := file.Write(itemXml); err != nil {
			return fmt.Errorf("error writing XML file: %s", err)
		}
		if _, err := file.Write([]byte("\n")); err != nil { // opcional: agregar un salto de línea entre objetos
			return fmt.Errorf("error writing new line to XML file: %s", err)
		}
	}

	return nil
}
