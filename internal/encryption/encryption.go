package encryption

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"errors"
)

// GenerateKeys generates an ECDSA key pair.
func GenerateKeys() (*ecdsa.PrivateKey, *ecdsa.PublicKey, error) {
	privKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		return nil, nil, err
	}
	return privKey, &privKey.PublicKey, nil
}

// GenerateProof generates and encrypts 1024 random strings of 1024 characters each using the public key.
func GenerateProof(pubKey *ecdsa.PublicKey) ([][]byte, [][]byte, error) {
	var encryptedStrings [][]byte
	var randomStrings [][]byte
	for i := 0; i < 1024; i++ {
		randomData := make([]byte, 1024)
		_, err := rand.Read(randomData)
		if err != nil {
			return nil, nil, err
		}

		randomStrings = append(randomStrings, randomData)
		nonce, ciphertext, err := EncryptMessage(randomData, pubKey)
		if err != nil {
			return nil, nil, err
		}

		// Łączymy nonce i ciphertext: [nonce || ciphertext]
		encrypted := append(nonce, ciphertext...)
		encryptedStrings = append(encryptedStrings, encrypted)
	}
	return encryptedStrings, randomStrings, nil
}

// VerifyProof verifies if the decrypted strings match the original random strings.
func VerifyProof(randomStrings [][]byte, decryptedStrings [][]byte) bool {
	if len(randomStrings) != len(decryptedStrings) {
		return false
	}
	for i := range randomStrings {
		if !compareByteSlices(randomStrings[i], decryptedStrings[i]) {
			return false
		}
	}
	return true
}

// compareByteSlices compares two byte slices for equality.
func compareByteSlices(a, b []byte) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

// EncryptMessage szyfruje wiadomość, używając wyłącznie klucza publicznego.
func EncryptMessage(message []byte, pubKey *ecdsa.PublicKey) ([]byte, []byte, error) {
	// Klucz do AES-GCM to SHA-256 z bajtowej reprezentacji klucza publicznego
	pubBytes := elliptic.Marshal(pubKey.Curve, pubKey.X, pubKey.Y)
	key := sha256.Sum256(pubBytes)

	block, err := aes.NewCipher(key[:])
	if err != nil {
		return nil, nil, err
	}
	aead, err := cipher.NewGCM(block)
	if err != nil {
		return nil, nil, err
	}

	nonce := make([]byte, aead.NonceSize())
	if _, err := rand.Read(nonce); err != nil {
		return nil, nil, err
	}

	ciphertext := aead.Seal(nil, nonce, message, nil)
	return nonce, ciphertext, nil
}

// DecryptMessage odszyfrowuje wiadomość, korzystając z klucza prywatnego
// (klient ma do niego dostęp, więc serwer tego nie robi).
func DecryptMessage(nonce, ciphertext []byte, privKey *ecdsa.PrivateKey) ([]byte, error) {
	pub := &privKey.PublicKey
	pubBytes := elliptic.Marshal(pub.Curve, pub.X, pub.Y)
	key := sha256.Sum256(pubBytes)

	block, err := aes.NewCipher(key[:])
	if err != nil {
		return nil, err
	}
	aead, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}
	if len(nonce) != aead.NonceSize() {
		return nil, errors.New("invalid nonce size")
	}

	plaintext, err := aead.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return nil, err
	}
	return plaintext, nil
}
