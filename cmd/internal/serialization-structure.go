package internal

import (
	"encoding/json"
	"fmt"
	"math/big"
)

func JSONStringToBigIntSlice(jsonString string) ([]*big.Int, error) {
	// Deserializuj JSON do []string
	var stringSlice []string
	err := json.Unmarshal([]byte(jsonString), &stringSlice)
	if err != nil {
		return nil, err
	}

	// Konwertuj []string na []*big.Int
	bigIntSlice := make([]*big.Int, len(stringSlice))
	for i, str := range stringSlice {
		bigInt := new(big.Int)
		_, ok := bigInt.SetString(str, 10) // Konwersja z bazy 10
		if !ok {
			return nil, fmt.Errorf("failed to convert string to big.Int: %s", str)
		}
		bigIntSlice[i] = bigInt
	}

	return bigIntSlice, nil
}
