package utils

import "sync"

func GoRoutineWorker(wg *sync.WaitGroup, jobs <-chan map[string]string, results chan<- []string, headers []string) {
	defer wg.Done()
	for job := range jobs {
		record := make([]string, len(headers))
		for i, header := range headers {
			value := job[header]
			record[i] = value
		}
		results <- record
	}
}