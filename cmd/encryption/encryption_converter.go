package encryption

import (
	"encoding/hex"
	"fmt"
	"math/big"

	"github.com/consensys/gnark-crypto/ecc/bn254"
)

func G1AffineToString(y bn254.G1Affine) string {
	bytes := y.Marshal()
	return hex.EncodeToString(bytes)
}

func StringToG1Affine(s string) (bn254.G1Affine, error) {
	var y bn254.G1Affine

	bytes, err := hex.DecodeString(s)
	if err != nil {
		return y, err
	}

	err = y.Unmarshal(bytes)
	if err != nil {
		return y, err
	}

	return y, nil
}

func PublicKeyStringToBigInt(publicKey string) (*big.Int, error) {

	v := new(big.Int)
	_, success := v.SetString(publicKey, 16)
	if !success {
		return nil, fmt.Errorf("failed to convert public key string to big int")
	}

	return v, nil
}
