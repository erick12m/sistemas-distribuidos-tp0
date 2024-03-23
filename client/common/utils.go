package common

import (
	"fmt"
	"os"
)

func serialize(data ClientData, agency string) string {
	return fmt.Sprintf(
		"%v,%v,%v,%v,%v,%v",
		agency,
		data.Name,
		data.LastName,
		data.Document,
		data.Birthdate,
		data.Number,
	)
}

func openBetFile(betsFile string) (*os.File, error) {
	file, err := os.Open(betsFile)
	if err != nil {
		return nil, err
	}
	return file, nil
}
