package main

import (
	"bufio"
	"encoding/csv"
	"io"
	"os"
)

func loadUserDataFromFile(filename string) map[string]string {
	fp, _ := os.Open(filename)
	return loadCSV(bufio.NewReader(fp))
}

func loadCSV(reader io.Reader) map[string]string {
	var data = map[string]string{}

	csvReader := csv.NewReader(reader)
	csvReader.Comma = ','
	for {
		record, err := csvReader.Read()
		if err == io.EOF {
			break
		}
		if len(record) > 1 {
			data[record[0]] = record[1]
		}
	}
	return data
}
