package main

import (
	"flag"
	"fmt"
	"github/file_converter/csv"
	"github/file_converter/xml"
	"time"
)

type ConvertFunc func(string) error

func main() {
	start := time.Now()

	path := flag.String("path", "", "Ruta al archivo a convertir")
	to := flag.String("to", "", "Formato de salida (json, csv, xml, yaml)")
	flag.Parse()

	fmt.Println("Convirtiendo archivo:", *path, "a formato:", *to)

	convertOptions := map[string]ConvertFunc{
		"json": csv.ConvertToCsv,
		"csv":  csv.ConvertToCsv,
		"xml":  xml.ConvertToXml,
		"yaml": csv.ConvertToCsv,
	}

	if *path == "" || *to == "" {
		fmt.Println("Error: Debes proporcionar la ruta al archivo y el formato de salida.")
		return
	}

	value := convertOptions[*to]

	if value == nil {
		fmt.Println("Error: seleccione una de estas opciones: json, csv, xml, yaml.")
		return
	}

	if convertFunc, ok := convertOptions[*to]; ok {
		err := convertFunc(*path)
		fmt.Println(err)
	}

	duration := time.Since(start)
	seconds := duration.Seconds()

	fmt.Printf("Conversión completada en %.2f segundos\n", seconds)
}
