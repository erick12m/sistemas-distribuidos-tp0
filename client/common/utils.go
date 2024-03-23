package common

import (
	"fmt"
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
