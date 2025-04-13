package internal

import (
	"math/big"
	"time"

	"github.com/consensys/gnark-crypto/ecc/bn254"
)

type SchnorrChallengeToSave struct {
	PublicKey string
	R         bn254.G1Affine
	E         *big.Int
	Expiry    time.Time
}

type SchnorrChallengeToSend struct {
	R *big.Int
	E *big.Int
}

type FeigeFiatShamirChallengeToSave struct {
	N      *big.Int
	E      *big.Int
	Expiry time.Time
}

type FeigeFiatShamirChallengeToSend struct {
	C1N        string
	EncryptedN string
	C1e        string
	EncryptedE string
}

type SigmaChallengeToSave struct {
	E      *big.Int
	T      bn254.G1Affine
	Expiry time.Time
}

type SigmaChallengeToSend struct {
	C1e        string
	EncryptedE string
	C1r        string
	EncryptedR string
}
