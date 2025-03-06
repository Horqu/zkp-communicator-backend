package encryption

import (
	"testing"
)

func TestEncryptionFlow(t *testing.T) {
	// Klient generuje parę kluczy ECDSA
	privKey, pubKey, err := GenerateKeys()
	if err != nil {
		t.Fatalf("Failed to generate keys: %v", err)
	}

	// Serwer generuje proof używając klucza publicznego klienta
	encryptedStrings, randomStrings, err := GenerateProof(pubKey)
	if err != nil {
		t.Fatalf("Failed to generate proof: %v", err)
	}

	// Klient odszyfrowuje proof używając klucza prywatnego
	var decryptedStrings [][]byte
	for _, encryptedString := range encryptedStrings {
		nonceSize := 12 // AES-GCM standard nonce size
		nonce := encryptedString[:nonceSize]
		ciphertext := encryptedString[nonceSize:]

		plaintext, err := DecryptMessage(nonce, ciphertext, privKey)
		if err != nil {
			t.Fatalf("Failed to decrypt message: %v", err)
		}
		decryptedStrings = append(decryptedStrings, plaintext)
	}

	// Serwer weryfikuje, czy odszyfrowane ciągi znaków są poprawne
	if !VerifyProof(randomStrings, decryptedStrings) {
		t.Fatalf("Proof verification failed")
	}
}
