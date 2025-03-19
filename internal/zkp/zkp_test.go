package zkp

import (
	"testing"
)

func TestZKPFlow(t *testing.T) {
	Gnark_crypto_main()
}

func TestSchnorrProof(t *testing.T) {
	Schnorr_proof()
	// Testowanie generowania i weryfikacji dowodu Schnorra
}

func TestEncryptionDecryption(t *testing.T) {
	encryption_decryption_test()
	// Testowanie szyfrowania i deszyfrowania
}
