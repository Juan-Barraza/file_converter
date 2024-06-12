package xml

import (
	"encoding/xml"
	"fmt"
	"os"
	"github.com/DeijoseDevelop/file_converter/converter"
)

func init() {
	converter.RegisterReadConvertFunc("xml", ReadXml)

}

func ReadXml(path string) ([]map[string]interface{}, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("error opening the file: %s", err)
	}
	defer file.Close()

	decoder := xml.NewDecoder(file)
	var records []map[string]interface{}
	var currentElement string
	var currentRecord map[string]interface{}

	for {
		token, err := decoder.Token()
		if err != nil {
			if err.Error() == "EOF" {
				break
			}
			return nil, fmt.Errorf("error decoding XML: %s", err)
		}

		switch element := token.(type) {
		case xml.StartElement:
			currentElement = element.Name.Local
			if currentElement == "record" {
				currentRecord = make(map[string]interface{})
			}
		case xml.EndElement:
			if element.Name.Local == "record" && currentRecord != nil {
				records = append(records, currentRecord)
				currentRecord = nil
			}
		case xml.CharData:
			if currentRecord != nil {
				currentRecord[currentElement] = string(element)
			}
		}
	}

	return records, nil
}
