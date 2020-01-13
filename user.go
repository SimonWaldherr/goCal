package main

import (
	"bufio"
	"encoding/csv"
	"io"
	"os"
	"strings"
)

type accessType string

const (
	readonly  accessType = "RO"
	writeonly accessType = "WO"
	readwrite accessType = "RW"
)

type userData struct {
	name     string
	password string
	access   accessType
}

func loadUserDataFromFile(filename string) map[string]userData {
	fp, _ := os.Open(filename)
	return loadCSV(bufio.NewReader(fp))
}

func loadCSV(reader io.Reader) map[string]userData {
	var data = map[string]userData{}
	var accT accessType

	csvReader := csv.NewReader(reader)
	csvReader.Comma = ','
	for {
		record, err := csvReader.Read()
		if err == io.EOF {
			break
		}
		if len(record) > 1 {
			if len(record) > 2 {
				switch strings.ToUpper(record[2]) {
				case "RO":
					accT = readonly
				case "WO":
					accT = writeonly
				case "RW":
					accT = readwrite
				default:
					accT = readonly
				}
			} else {
				accT = readonly
			}
			data[record[0]] = userData{
				name:     record[0],
				password: record[1],
				access:   accT,
			}
		}
	}
	return data
}
