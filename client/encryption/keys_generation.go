package encryption

import (
	"crypto/rand"
	"encoding/hex"
	"math/big"

	"github.com/consensys/gnark-crypto/ecc/bn254"
)

func GeneratePrivateKey() big.Int {
	x, err := rand.Int(rand.Reader, bn254.ID.ScalarField())
	if err != nil {
		panic(err)
	}
	return *x
}

func GeneratePublicKey(x big.Int) bn254.G1Affine {
	var g bn254.G1Affine
	g.X.SetString("1")
	g.Y.SetString("2")

	if !g.IsOnCurve() {
		panic("Point is not on curve")
	}

	var y bn254.G1Affine
	y.ScalarMultiplication(&g, &x)

	return y
}

func PublicKeyToString(y bn254.G1Affine) string {
	bytes := y.Marshal()
	return hex.EncodeToString(bytes)
}

func StringToPublicKey(s string) (bn254.G1Affine, error) {
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
