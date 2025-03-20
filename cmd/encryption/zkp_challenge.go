package encryption

import (
	"crypto/rand"
	"fmt"

	"github.com/consensys/gnark-crypto/ecc/bn254"
)

func GenerateZkSnarkProof() {
}

func VerifyZkSnarkProof(check bn254.G1Affine, y bn254.G1Affine) bool {

	if check.Equal(&y) {
		return true
	} else {
		return false
	}

}

func Gnark_crypto_main() {
	// Generator punktu na krzywej BN254
	var g bn254.G1Affine
	g.X.SetString("1") // Wartość generatora X
	g.Y.SetString("2") // Wartość generatora Y

	if !g.IsOnCurve() {
		panic("Punkt nie należy do krzywej")
	}
	// Wygenerowanie sekretnego klucza x
	// x, err := rand.Int(rand.Reader, fr.Modulus())
	x, err := rand.Int(rand.Reader, bn254.ID.ScalarField())
	if err != nil {
		panic(err)
	}

	fmt.Println("Sekretny klucz x:", x)

	// Obliczenie y = g^x
	var y bn254.G1Affine
	y.ScalarMultiplication(&g, x)

	fmt.Println("Publiczny klucz y:", y.String())

	// Weryfikacja dowodu (sprawdzenie, czy x rzeczywiście generuje y)
	var check bn254.G1Affine
	check.ScalarMultiplication(&g, x) // to musi zwrocic klient do serwera

	if check.Equal(&y) {
		fmt.Println("Dowód poprawny: posiadacz zna x")
	} else {
		fmt.Println("Dowód błędny: x nie pasuje")
	}
}
