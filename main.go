package main

import (
	"flag"
	"fmt"
	"path/filepath"
	"runtime"
	"time"

	"github.com/DeijoseDevelop/file_converter/converter"
	"github.com/DeijoseDevelop/file_converter/json"
	"github.com/DeijoseDevelop/file_converter/csv"
	"github.com/DeijoseDevelop/file_converter/xml"
	"github.com/DeijoseDevelop/file_converter/yaml"
)

type ConvertFunc func(string, string) error

func main() {
	start := time.Now()

	path := flag.String("path", "", "Ruta al archivo a convertir")
	to := flag.String("to", "", "Formato de salida (json, csv, xml, yaml)")
	flag.Parse()

	var memStats runtime.MemStats
    runtime.ReadMemStats(&memStats)
    initialAlloc := memStats.Alloc

	fmt.Println("Convirtiendo archivo:", *path, "a formato:", *to)

	convertOptions := map[string]ConvertFunc{
		"json": json.ConvertToJson,
		"csv":  csv.ConvertToCsv,
		"xml":  xml.ConvertToXml,
		"yaml": yaml.ConvertToYaml,
	}

	converter.RegisterReadConvertFunc("json", json.ReadJson)
	converter.RegisterReadConvertFunc("csv", csv.ReadCSV)
	converter.RegisterReadConvertFunc("xml", xml.ReadXml)
	converter.RegisterReadConvertFunc("yaml", yaml.ReadYaml)

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
		err := convertFunc(*path, filepath.Ext(*path)[1:])
		fmt.Printf("Error: %s\n", err)
	}

	runtime.ReadMemStats(&memStats)
    finalAlloc := memStats.Alloc

	duration := time.Since(start)
	seconds := float64(duration.Nanoseconds()) / 1000000000.0
    cpuUsagePercent := 100 * seconds / (seconds * float64(runtime.NumCPU()))
    memoryUsageMB := float64(finalAlloc - initialAlloc) / 1024.0 / 1024.0

	fmt.Printf("Conversión completada en %.2f segundos\n", seconds)
    fmt.Printf("Uso de memoria: %.2f MB\n", memoryUsageMB)
    fmt.Printf("Uso de CPU: %.2f%%\n", cpuUsagePercent)
}
