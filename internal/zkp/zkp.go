package zkp

import (
	"math/big"
	"net/http"

	"github.com/gin-gonic/gin"

	"crypto/rand"
	"crypto/sha256"
	"fmt"

	"github.com/consensys/gnark-crypto/ecc/bn254"
	"github.com/consensys/gnark-crypto/ecc/bn254/fr"
)

// GenerateProofHandler generates ZKP proof
func GenerateProofHandler(c *gin.Context) {
	// Logic for generating ZKP proof
	// TODO: Implementation ZKP proof generation

	c.JSON(http.StatusOK, gin.H{
		"proof": "sample_proof_generated",
	})
}

// VerifyProofHandler verifies ZKP proof
func VerifyProofHandler(c *gin.Context) {
	// Logic for verifying ZKP proof
	// TODO: Implementation of ZKP proof verification

	c.JSON(http.StatusOK, gin.H{
		"valid": true,
	})
}

// ZK-SNARK (Zero-Knowledge Succinct Non-Interactive Argument of Knowledge)

func Gnark_crypto_main() {
	// Generator punktu na krzywej BN254
	var g bn254.G1Affine
	g.X.SetString("11")  // Wartość generatora X
	g.Y.SetString("222") // Wartość generatora Y

	// Wygenerowanie sekretnego klucza x
	x, err := rand.Int(rand.Reader, fr.Modulus())
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

/*
Serwer musi przekazać klientowi ustalony punkt bazowy (generator) – np. g – oraz ewentualnie
losowe wyzwanie (lub sposób jego obliczenia).
Klient odsyła serwerowi swój punkt publiczny y, zobowiązanie R i odpowiedź s.
Serwer weryfikuje, czy dowód jest poprawny,
nawet jeśli nie ma wcześniej zarejestrowanego klucza publicznego.
*/
func Schnorr_proof() {
	// Generator na krzywej BN254
	var g bn254.G1Affine
	g.X.SetString("1")
	g.Y.SetString("2")

	// Generowanie klucza prywatnego x
	x, err := rand.Int(rand.Reader, bn254.ID.ScalarField())
	fmt.Println("Klucz prywatny x:", x)
	if err != nil {
		panic(err)
	}

	// Obliczenie klucza publicznego y = g^x
	var y bn254.G1Affine
	y.ScalarMultiplication(&g, x)
	fmt.Println("Klucz publiczny y:", y.String())

	// Wygenerowanie losowej wartości r
	r, err := rand.Int(rand.Reader, bn254.ID.ScalarField())
	if err != nil {
		panic(err)
	}

	// Obliczenie zobowiązania R = g^r
	var R bn254.G1Affine
	R.ScalarMultiplication(&g, r)

	// Obliczenie wyzwania e = H(R || y)
	h := sha256.New()
	h.Write(R.Marshal())
	h.Write(y.Marshal())
	e := new(big.Int).SetBytes(h.Sum(nil))
	e.Mod(e, bn254.ID.ScalarField())

	// Obliczenie odpowiedzi s = r + e * x
	s := new(big.Int).Mul(e, x)
	s.Add(s, r)
	s.Mod(s, bn254.ID.ScalarField())

	// Weryfikacja dowodu: sprawdzamy, czy g^s = R + y^e
	var gS, yE, RplusYE bn254.G1Affine
	gS.ScalarMultiplication(&g, s)
	yE.ScalarMultiplication(&y, e)
	RplusYE.Add(&R, &yE)

	if gS.Equal(&RplusYE) {
		fmt.Println("Dowód Schnorra poprawny: klient zna klucz prywatny.")
	} else {
		fmt.Println("Dowód Schnorra niepoprawny: klient nie zna klucza.")
	}
}
