package encryption

import (
	"crypto/rand"
	"log"
	"testing"

	"github.com/consensys/gnark-crypto/ecc/bn254"
)

func TestSchnorrProof(t *testing.T) {
	// Generator na krzywej BN254
	var g bn254.G1Affine
	g.X.SetString("1")
	g.Y.SetString("2")

	// Generowanie klucza prywatnego x
	x, err := rand.Int(rand.Reader, bn254.ID.ScalarField())
	if err != nil {
		t.Fatalf("Błąd generowania klucza prywatnego: %v", err)
	}

	// Obliczenie klucza publicznego y = g^x
	var y bn254.G1Affine
	y.ScalarMultiplication(&g, x)

	yString := G1AffineToString(y)

	// Generowanie wyzwania i zobowiązania
	e, r := GenerateSchnorrChallenge(yString)
	s := GenerateSchnorrProof(yString, x, e, r)

	// Weryfikacja dowodu
	var R bn254.G1Affine
	R.ScalarMultiplication(&g, r)
	if !VerifySchnorrProof(R, e, s, yString) {
		t.Errorf("Dowód Schnorra nie został zweryfikowany poprawnie")
	} else {
		t.Log("Dowód Schnorra zweryfikowany poprawnie")
	}
}

func TestSchnorrProofFailure(t *testing.T) {
	// Generator na krzywej BN254
	var g bn254.G1Affine
	g.X.SetString("1")
	g.Y.SetString("2")

	// Generowanie klucza prywatnego x
	x, err := rand.Int(rand.Reader, bn254.ID.ScalarField())
	if err != nil {
		t.Fatalf("Błąd generowania klucza prywatnego: %v", err)
	}

	// Obliczenie klucza publicznego y = g^x
	var y bn254.G1Affine
	y.ScalarMultiplication(&g, x)

	yString := G1AffineToString(y)

	// Generowanie wyzwania i zobowiązania
	e, r := GenerateSchnorrChallenge(yString)
	s := GenerateSchnorrProof(yString, x, e, r)

	// Celowe wprowadzenie błędu: zmodyfikowanie klucza publicznego
	var yModified bn254.G1Affine
	yModified.X.SetString("123456789") // Nieprawidłowy punkt na krzywej
	yModified.Y.SetString("987654321")
	yModifiedString := G1AffineToString(yModified)

	// Weryfikacja dowodu z błędnym kluczem publicznym
	var R bn254.G1Affine
	R.ScalarMultiplication(&g, r)
	if VerifySchnorrProof(R, e, s, yModifiedString) {
		t.Errorf("Dowód Schnorra został zweryfikowany poprawnie, mimo że powinien zakończyć się niepowodzeniem")
	} else {
		t.Log("Dowód Schnorra poprawnie nie został zweryfikowany (test zakończony niepowodzeniem zgodnie z oczekiwaniami)")
	}
}

func TestFeigeFiatShamirProof(t *testing.T) {

	var g bn254.G1Affine
	g.X.SetString("1")
	g.Y.SetString("2")

	privateKey, err := rand.Int(rand.Reader, bn254.ID.ScalarField())
	if err != nil {
		t.Fatalf("Błąd generowania klucza prywatnego: %v", err)
	}

	N, e := GenerateFeigeFiatShamirChallenge()

	log.Println("Klucz prywatny:", privateKey.String())

	x, y, v := GenerateFeigeFiatShamirProof(privateKey, N, e)

	if VerifyFeigeFiatShamir(x, y, v, N, e) {
		t.Log("Dowód Feige-Fiat-Shamir zweryfikowany poprawnie")
	} else {
		t.Errorf("Dowód Feige-Fiat-Shamir nie został zweryfikowany poprawnie")
	}

}
