package xml

import (
	"encoding/xml"
	"fmt"
	"os"
)

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
	var inRecord bool

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
			if element.Name.Local == "record" {
				currentRecord = make(map[string]interface{})
				inRecord = true
			} else if inRecord {
				currentElement = element.Name.Local
			}
		case xml.EndElement:
			if element.Name.Local == "record" && currentRecord != nil {
				records = append(records, currentRecord)
				currentRecord = nil
				inRecord = false
			}
			currentElement = ""
		case xml.CharData:
			if inRecord && currentElement != "" {
				trimmedValue := string(element)
				if len(trimmedValue) > 0 {
					currentRecord[currentElement] = trimmedValue
				}
			}
		}
	}

	return records, nil
}
