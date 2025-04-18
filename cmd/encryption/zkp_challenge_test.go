package encryption

import (
	"crypto/rand"
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

	t.Logf("Klucz prywatny: %s", x.String())

	// Obliczenie klucza publicznego y = g^x
	var y bn254.G1Affine
	y.ScalarMultiplication(&g, x)

	yString := G1AffineToString(y)

	t.Logf("Klucz publiczny: %s", yString)

	// Generowanie wyzwania i zobowiązania
	e, r, R := GenerateSchnorrChallenge(yString)
	s := GenerateSchnorrProof(x, e, r)

	// // Weryfikacja dowodu
	if !VerifySchnorrProof(R, e, s, yString) {
		t.Errorf("Dowód Schnorra nie został zweryfikowany poprawnie")
	} else {
		t.Log("Dowód Schnorra zweryfikowany poprawnie")
	}
}

func TestSchnorrProofSpecific(t *testing.T) {
	// Generator na krzywej BN254
	var g bn254.G1Affine
	g.X.SetString("1")
	g.Y.SetString("2")

	// Generowanie klucza prywatnego x
	x := "10099027144552925486150250931935288803469455528408931520761903588766227001874"
	xBigInt, err := PublicKeyStringToBigInt(x)
	if err != nil {
		t.Fatalf("Błąd konwersji klucza prywatnego: %v", err)
	}

	var y bn254.G1Affine
	y.ScalarMultiplication(&g, xBigInt)
	t.Logf("Klucz publiczny: %s", G1AffineToString(y))

	yString := "241f6a1a8e0acb84e3546466cc81a04d0d6fcd6b0319f6aa44467d77f72224f219ab9482079ff001251e582e90b4d287c75244a3a4ccae5b91be497e8e580f78"
	// yString := G1AffineToString(y)

	// Generowanie wyzwania i zobowiązania
	e, r, R := GenerateSchnorrChallenge(yString)
	s := GenerateSchnorrProof(xBigInt, e, r)

	t.Logf("Wyzwanie e: %s", e.String())
	t.Logf("Zobowiązanie r: %s", r.String())
	t.Logf("Zobowiązanie R: %s", R.String())
	t.Logf("Dowód s: %s", s.String())

	// Weryfikacja dowodu
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
	e, r, _ := GenerateSchnorrChallenge(yString)
	s := GenerateSchnorrProof(x, e, r)

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

	x, y, v := GenerateFeigeFiatShamirProof(privateKey, N, e)

	if VerifyFeigeFiatShamir(x, y, v, N, e) {
		t.Log("Dowód Feige-Fiat-Shamir zweryfikowany poprawnie")
	} else {
		t.Errorf("Dowód Feige-Fiat-Shamir nie został zweryfikowany poprawnie")
	}

}

func TestSigmaProof(test *testing.T) {
	// Generator na krzywej BN254
	var g bn254.G1Affine
	g.X.SetString("1")
	g.Y.SetString("2")

	// Generowanie klucza prywatnego
	privateKey, err := rand.Int(rand.Reader, bn254.ID.ScalarField())
	if err != nil {
		test.Fatalf("Błąd generowania klucza prywatnego: %v", err)
	}

	// Obliczenie klucza publicznego publicKey = g^privateKey
	var publicKey bn254.G1Affine
	publicKey.ScalarMultiplication(&g, privateKey)

	// Generowanie wyzwania i zobowiązania
	e, r, _ := GenerateSigmaChallenge()

	// Obliczenie zobowiązania t = g^r
	var t bn254.G1Affine
	t.ScalarMultiplication(&g, r)

	// Generowanie dowodu
	s := GenerateSigmaProof(privateKey, e, r)

	// Weryfikacja dowodu
	if VerifySigmaProof(&t, e, s, publicKey) {
		test.Log("Dowód Sigma Protocol zweryfikowany poprawnie")
	} else {
		test.Errorf("Dowód Sigma Protocol nie został zweryfikowany poprawnie")
	}
}

func TestSigmaProofFailure(test *testing.T) {
	// Generator na krzywej BN254
	var g bn254.G1Affine
	g.X.SetString("1")
	g.Y.SetString("2")

	// Generowanie klucza prywatnego
	privateKey, err := rand.Int(rand.Reader, bn254.ID.ScalarField())
	if err != nil {
		test.Fatalf("Błąd generowania klucza prywatnego: %v", err)
	}

	// Obliczenie klucza publicznego publicKey = g^privateKey
	var publicKey bn254.G1Affine
	publicKey.ScalarMultiplication(&g, privateKey)

	// Generowanie wyzwania i zobowiązania
	e, r, _ := GenerateSigmaChallenge()

	// Obliczenie zobowiązania t = g^r
	var t bn254.G1Affine
	t.ScalarMultiplication(&g, r)

	// Generowanie dowodu
	s := GenerateSigmaProof(privateKey, e, r)

	// Celowe wprowadzenie błędu: zmodyfikowanie klucza publicznego
	var modifiedPublicKey bn254.G1Affine
	modifiedPublicKey.X.SetString("123456789") // Nieprawidłowy punkt na krzywej
	modifiedPublicKey.Y.SetString("987654321")

	// Weryfikacja dowodu z błędnym kluczem publicznym
	if VerifySigmaProof(&t, e, s, modifiedPublicKey) {
		test.Errorf("Dowód Sigma Protocol został zweryfikowany poprawnie, mimo że powinien zakończyć się niepowodzeniem")
	} else {
		test.Log("Dowód Sigma Protocol poprawnie nie został zweryfikowany (test zakończony niepowodzeniem zgodnie z oczekiwaniami)")
	}
}
