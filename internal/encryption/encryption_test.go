package encryption

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"go.mau.fi/libsignal/ecc"
	"go.mau.fi/libsignal/keys/identity"
)

func TestGenerateKeyPair(t *testing.T) {
	keyPair, err := GenerateKeyPair()
	assert.NoError(t, err)
	assert.NotNil(t, keyPair)
	assert.NotNil(t, keyPair.PublicKey())
	assert.NotNil(t, keyPair.PrivateKey())
}

func TestEncryptDecryptMessage(t *testing.T) {
	keyPair, err := GenerateKeyPair()
	assert.NoError(t, err)

	senderKeyPair, err := GenerateKeyPair()
	assert.NoError(t, err)
	senderIdentityKey := identity.NewKey(senderKeyPair.PublicKey())

	receiverKeyPair, err := GenerateKeyPair()
	assert.NoError(t, err)
	receiverIdentityKey := identity.NewKey(receiverKeyPair.PublicKey())

	message := "Hello, World!"
	publicKey, ok := keyPair.PublicKey().(*ecc.DjbECPublicKey)
	assert.True(t, ok)

	ciphertext, err := EncryptMessage(message, publicKey, senderIdentityKey, receiverIdentityKey)
	assert.NoError(t, err)
	assert.NotNil(t, ciphertext)

	privateKey, ok := keyPair.PrivateKey().(*ecc.DjbECPrivateKey)
	assert.True(t, ok)

	decryptedMessage, err := DecryptMessage(ciphertext, privateKey)
	assert.NoError(t, err)
	assert.Equal(t, message, decryptedMessage)
}

func TestGenerateGroupKey(t *testing.T) {
	groupKey, err := GenerateGroupKey()
	assert.NoError(t, err)
	assert.NotNil(t, groupKey)
	assert.NotNil(t, groupKey.PublicKey())
	assert.NotNil(t, groupKey.PrivateKey())
}

func TestGenerateNonce(t *testing.T) {
	nonce, err := GenerateNonce(12)
	assert.NoError(t, err)
	assert.NotNil(t, nonce)
	assert.Equal(t, 12, len(nonce))
}

func TestGenerateRandomBytes(t *testing.T) {
	randomBytes, err := GenerateRandomBytes(16)
	assert.NoError(t, err)
	assert.NotNil(t, randomBytes)
	assert.Equal(t, 16, len(randomBytes))
}
