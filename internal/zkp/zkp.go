package zkp

import (
	"encoding/hex"
	"io"
	"math/big"
	"net/http"

	"github.com/gin-gonic/gin"

	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"fmt"

	"github.com/consensys/gnark-crypto/ecc/bn254"
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

// encryptionDecryptionTest pokazuje przykład szyfrowania i deszyfrowania tekstu
func encryption_decryption_test() {
	// Generator na krzywej BN254
	var g bn254.G1Affine
	g.X.SetString("1")
	g.Y.SetString("2")

	// Klucz prywatny odbiorcy (x)
	x := new(big.Int)
	x.SetString("12065182636985977436074843367923784964065411516963092648656448500555354703849", 10)

	// Klucz publiczny odbiorcy (y = g^x)
	var y bn254.G1Affine
	y.ScalarMultiplication(&g, x)

	fmt.Println("Publiczny klucz y:", y.String())

	// Dowolny tekst do zaszyfrowania
	plaintext := "To jest tajna wiadomość."

	// Szyfrowanie (zwraca punkt C1 oraz zaszyfrowaną treść)
	C1, cipherBytes := encryptText(plaintext, &g, &y)

	// Wyświetlamy efekt szyfrowania
	fmt.Println("C1 (punkt na krzywej):", C1)
	fmt.Println("Zaszyfrowana wiadomość (AES):", hex.EncodeToString(cipherBytes))

	// Odszyfrowanie
	decryptedText := decryptText(C1, cipherBytes, x)
	fmt.Println("Odszyfrowana wiadomość:", decryptedText)
}

// encryptText – szyfruje tekst przy pomocy klucza publicznego y
// Zwraca:
//  1. C1 = g^k (bn254.G1Affine) – punkt pomocniczy
//  2. zaszyfrowany tekst ([]byte), zaszyfrowany AES-em
func encryptText(plaintext string, g, y *bn254.G1Affine) (bn254.G1Affine, []byte) {
	// 1. Wygeneruj losową wartość k
	k, err := rand.Int(rand.Reader, bn254.ID.ScalarField())
	if err != nil {
		panic(err)
	}

	// 2. Oblicz C1 = g^k
	var C1 bn254.G1Affine
	C1.ScalarMultiplication(g, k)

	// 3. Oblicz punkt shared = y^k
	var shared bn254.G1Affine
	shared.ScalarMultiplication(y, k)

	// 4. Z punktu shared twórz klucz symetryczny (np. AES) przez zhashowanie
	//    Marshal() serializuje punkt do bajtów; SHA-256 daje 32-bajtowy klucz
	key := sha256.Sum256(shared.Marshal())

	// 5. Szyfruj plaintext algorytmem AES (CFB szyfruje strumieniowo)
	block, err := aes.NewCipher(key[:])
	if err != nil {
		panic(err)
	}
	cipherBytes := make([]byte, aes.BlockSize+len(plaintext))
	iv := cipherBytes[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		panic(err)
	}

	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(cipherBytes[aes.BlockSize:], []byte(plaintext))

	return C1, cipherBytes
}

// decryptText – deszyfruje tekst przy pomocy klucza prywatnego x
// 1) Odtwarza ten sam klucz symetryczny z C1^x
// 2) Deszyfruje za pomocą AES
func decryptText(C1 bn254.G1Affine, cipherBytes []byte, x *big.Int) string {
	// 1. Oblicz punkt shared = C1^x
	var shared bn254.G1Affine
	shared.ScalarMultiplication(&C1, x)

	// 2. Hashujemy punkt shared, by uzyskać klucz AES
	key := sha256.Sum256(shared.Marshal())

	// 3. Deszyfrujemy dane
	block, err := aes.NewCipher(key[:])
	if err != nil {
		panic(err)
	}

	iv := cipherBytes[:aes.BlockSize]
	msg := cipherBytes[aes.BlockSize:]

	stream := cipher.NewCFBDecrypter(block, iv)
	stream.XORKeyStream(msg, msg)

	return string(msg)
}
