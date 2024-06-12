package converter



type ReadConvertFunc func(string) ([]map[string]any, error)

var readConvertOptions = map[string]ReadConvertFunc{}

func RegisterReadConvertFunc(format string, fn ReadConvertFunc) {
	readConvertOptions[format] = fn
}

func GetReadConvertFunc(format string) (ReadConvertFunc, bool) {
	fn, ok := readConvertOptions[format]
	return fn, ok
}

