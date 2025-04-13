package internal

import (
	"encoding/json"
	"math/big"
)

type SimplifiedContact struct {
	Username string `json:"username"`
	Status   string `json:"status"`
}

func BigIntSliceToJSONString(slice []*big.Int) (string, error) {
	// Konwertuj []big.Int na []string
	stringSlice := make([]string, len(slice))
	for i, num := range slice {
		stringSlice[i] = num.String()
	}

	// Serializuj []string do JSON
	jsonBytes, err := json.Marshal(stringSlice)
	if err != nil {
		return "", err
	}

	// Zwróć jako string
	return string(jsonBytes), nil
}
