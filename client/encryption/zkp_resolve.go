package encryption

import (
	"crypto/rand"
	"math/big"

	"github.com/consensys/gnark-crypto/ecc/bn254"
)

func GenerateSigmaProof(privateKey *big.Int, e *big.Int, r *big.Int) *big.Int {

	s := new(big.Int).Mul(e, privateKey)
	s.Add(s, r)
	s.Mod(s, bn254.ID.ScalarField())

	return s // received e, r | send s
}

func GenerateSchnorrProof(privateKey *big.Int, e *big.Int, r *big.Int) *big.Int {

	s := new(big.Int).Mul(e, privateKey)
	s.Add(s, r)
	s.Mod(s, bn254.ID.ScalarField())
	return s
}

func GenerateFeigeFiatShamirProof(privateKey *big.Int, N *big.Int, e *big.Int) ([]*big.Int, []*big.Int, *big.Int) {
	v := new(big.Int).Exp(privateKey, big.NewInt(2), N) // v = privateKey^2 mod N
	var xList []*big.Int
	var yList []*big.Int

	eCopy := new(big.Int).Set(e) // Kopia e

	for i := 0; i < 1024; i++ {
		// Wybór losowej wartości r
		r, err := rand.Int(rand.Reader, N)
		if err != nil {
			panic(err)
		}

		// Obliczenie x = r^2 mod N
		x := new(big.Int).Exp(r, big.NewInt(2), N)
		xList = append(xList, x)

		// Pobranie i-tego bitu z e
		bit := new(big.Int).And(eCopy, big.NewInt(1)) // Pobranie najmłodszego bitu
		eCopy.Rsh(eCopy, 1)                           // Przesunięcie w prawo o 1 bit

		// Obliczenie odpowiedzi y
		var y *big.Int
		if bit.Cmp(big.NewInt(0)) == 0 {
			y = r // Jeśli bit = 0, odpowiedź to r
		} else {
			y = new(big.Int).Mul(r, privateKey) // Jeśli bit = 1, odpowiedź to r * privateKey mod N
			y.Mod(y, N)
		}
		yList = append(yList, y)
	}

	return xList, yList, v
}
