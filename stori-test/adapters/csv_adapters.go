package adapters

import (
	"encoding/csv"
	"fmt"
	"os"
)

func ReadAllCSV(csvDir string) (*[][]string, error) {

	csvFile, err := os.Open(csvDir)

	if err != nil {
		return nil, fmt.Errorf("%s", err.Error())
	}

	defer csvFile.Close()

	reader := csv.NewReader(csvFile)
	csvLines, err := reader.ReadAll()

	if err != nil {
		return nil, fmt.Errorf("%s", err.Error())
	}

	return &csvLines, nil

}
