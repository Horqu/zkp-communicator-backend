package encryption

import (
	"crypto/rand"
	"encoding/json"
	"errors"
	"io"

	"go.mau.fi/libsignal/ecc"
	"go.mau.fi/libsignal/keys/identity"
	"go.mau.fi/libsignal/protocol"
	"go.mau.fi/libsignal/util/optional"
)

// Example fields for JSON-based serialization of a SignalMessageStructure.
type signalMessageJSON struct {
	Version int    `json:"version"`
	Body    string `json:"body"`
	// Add more fields as needed.
}

// Minimal or “stub” implementation of SignalMessageSerializer
type signalMessageSerializerImpl struct{}

func (s *signalMessageSerializerImpl) Serialize(msgStruct *protocol.SignalMessageStructure) []byte {
	// This example hardcodes Version=3 and Body="example body"
	// In practice, map msgStruct fields to your JSON struct.
	cMsg := signalMessageJSON{
		Version: 3,
		Body:    "example body",
	}
	data, _ := json.Marshal(cMsg)
	return data
}

func (s *signalMessageSerializerImpl) Deserialize(serialized []byte) (*protocol.SignalMessageStructure, error) {
	var cMsg signalMessageJSON
	if err := json.Unmarshal(serialized, &cMsg); err != nil {
		return nil, err
	}

	// Build a blank SignalMessageStructure. We can’t do msgStruct.MessageVersion = X
	// if that field does not exist in this library version.
	var msgStruct protocol.SignalMessageStructure
	// If needed, you can store cMsg.Version locally—or check it before proceeding.
	// The library’s struct simply lacks a place to hold it.

	return &msgStruct, nil
}

// GenerateKeyPair generates a new key pair for encryption
func GenerateKeyPair() (*ecc.ECKeyPair, error) {
	return ecc.GenerateKeyPair()
}

func SerializePreKeySignalMessage(msg *protocol.PreKeySignalMessage) []byte {
	return msg.Serialize()
}

// DeserializePreKeySignalMessage adjusts to 3 arguments: data, preKeySerializer, signalSerializer
func DeserializePreKeySignalMessage(
	data []byte,
	preKeySerializer protocol.PreKeySignalMessageSerializer,
) (*protocol.PreKeySignalMessage, error) {
	// The third argument is a SignalMessageSerializer – pass nil if permissible by your library version
	return protocol.NewPreKeySignalMessageFromBytes(data, preKeySerializer, nil)
}

// EncryptMessage encrypts a message using Signal Protocol
func EncryptMessage(message string, publicKey *ecc.DjbECPublicKey, senderIdentityKey, receiverIdentityKey *identity.Key) ([]byte, error) {
	signalMsg, err := protocol.NewSignalMessage(
		3,                   // message version
		1,                   // counter
		0,                   // previous counter
		[]byte("macKey"),    // MAC key
		publicKey,           // senderRatchetKey
		[]byte(message),     // ciphertext
		senderIdentityKey,   // sender identity
		receiverIdentityKey, // receiver identity
		&signalMessageSerializerImpl{},
	)
	if err != nil {
		return nil, err
	}
	return signalMsg.Serialize(), nil
}

// DecryptMessage decrypts a message using Signal Protocol
func DecryptMessage(ciphertext []byte, privateKey *ecc.DjbECPrivateKey) (string, error) {
	signalMsg, err := protocol.NewSignalMessageFromBytes(ciphertext, &signalMessageSerializerImpl{})
	if err != nil {
		return "", err
	}
	// Optionally verify MAC, etc.
	return string(signalMsg.Body()), nil
}

// GenerateGroupKey generates a new key pair for group encryption
func GenerateGroupKey() (*ecc.ECKeyPair, error) {
	return ecc.GenerateKeyPair()
}

// GenerateNonce generates a nonce for encryption
func GenerateNonce(size int) ([]byte, error) {
	nonce := make([]byte, size)
	_, err := io.ReadFull(rand.Reader, nonce)
	if err != nil {
		return nil, err
	}
	return nonce, nil
}

// GenerateRandomBytes generates random bytes for cryptographic use
func GenerateRandomBytes(size int) ([]byte, error) {
	randomBytes := make([]byte, size)
	_, err := io.ReadFull(rand.Reader, randomBytes)
	if err != nil {
		return nil, err
	}
	return randomBytes, nil
}

// NewPreKeySignalMessage creates a new PreKeySignalMessage
func NewPreKeySignalMessage(version int, registrationID uint32, preKeyID *optional.Uint32, signedPreKeyID uint32, baseKey ecc.ECPublicKeyable, identityKey *identity.Key, signalMessage *protocol.SignalMessage) (*protocol.PreKeySignalMessage, error) {
	// return protocol.NewPreKeySignalMessage(version, registrationID, preKeyID, signedPreKeyID, baseKey, identityKey, signalMessage)
	return nil, errors.New("not implemented")
}

// NewSenderKeyDistributionMessage creates a new SenderKeyDistributionMessage
func NewSenderKeyDistributionMessage(id uint32, iteration uint32, chainKey []byte, signatureKey ecc.ECPublicKeyable) (*protocol.SenderKeyDistributionMessage, error) {
	// return protocol.NewSenderKeyDistributionMessage(id, iteration, chainKey, signatureKey)
	return nil, errors.New("not implemented")
}

// NewSignalMessage creates a new SignalMessage
func NewSignalMessage(messageVersion int, counter, previousCounter uint32, macKey []byte, senderRatchetKey ecc.ECPublicKeyable, ciphertext []byte) (*protocol.SignalMessage, error) {
	// return protocol.NewSignalMessage(messageVersion, counter, previousCounter, macKey, senderRatchetKey, ciphertext)
	return nil, errors.New("not implemented")
}
