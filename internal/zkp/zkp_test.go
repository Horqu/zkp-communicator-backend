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
