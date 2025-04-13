package encryption

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"io"
	"math/big"

	"github.com/consensys/gnark-crypto/ecc/bn254"
)

var UserPrivateKey string

// encryptText – szyfruje tekst przy pomocy klucza publicznego y
// Zwraca:
//  1. C1 = g^k (bn254.G1Affine) – punkt pomocniczy
//  2. zaszyfrowany tekst ([]byte), zaszyfrowany AES-em
func EncryptText(plaintext string, y *bn254.G1Affine) (bn254.G1Affine, string) {

	var g bn254.G1Affine
	g.X.SetString("1")
	g.Y.SetString("2")

	// 1. Wygeneruj losową wartość k
	k, err := rand.Int(rand.Reader, bn254.ID.ScalarField())
	if err != nil {
		panic(err)
	}

	// 2. Oblicz C1 = g^k
	var C1 bn254.G1Affine
	C1.ScalarMultiplication(&g, k)

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

	cipherTextBase64 := base64.StdEncoding.EncodeToString(cipherBytes)

	return C1, cipherTextBase64
}

// decryptText – deszyfruje tekst przy pomocy klucza prywatnego x
// 1) Odtwarza ten sam klucz symetryczny z C1^x
// 2) Deszyfruje za pomocą AES
func DecryptText(C1 bn254.G1Affine, cipherTextBase64 string, x *big.Int) string {
	cipherBytes, err := base64.StdEncoding.DecodeString(cipherTextBase64)
	if err != nil {
		panic(err)
	}
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
