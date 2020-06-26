package quorum

import (
	"crypto/rand"
	"encoding/json"
	"testing"

	"golang.org/x/crypto/argon2"

	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/nacl/box"
)

func TestEncodePrivateKey_whenTypical(t *testing.T) {
	pub, priv, err := box.GenerateKey(rand.Reader)
	if err != nil {
		t.Fatal(err)
	}
	jsonString, _, err := toKeyDataJSON("", nil, priv[:], toStandardBase64EncodedString(pub[:]))
	assert.NoError(t, err)
	t.Log(toStandardBase64EncodedString(pub[:]))
	t.Log(jsonString)
}

func TestEncodePrivateKey_whenUsingPassword(t *testing.T) {
	pub, priv, err := box.GenerateKey(rand.Reader)
	if err != nil {
		t.Fatal(err)
	}
	password := "1234"
	jsonString, _, err := toKeyDataJSON(password, nil, priv[:], toStandardBase64EncodedString(pub[:]))
	assert.NoError(t, err)
	t.Log(toStandardBase64EncodedString(pub[:]))
	t.Log(jsonString)

	// let's try to reconstruct private key
	var data keyData
	assert.NoError(t, json.Unmarshal([]byte(jsonString), &data))
	keyInfo := data.Config
	aOpts := keyInfo.Data.ArgonOpts
	salt, err := keyInfo.Data.ArgonSalt.bytes()
	assert.NoError(t, err)
	var hash []byte
	switch aOpts.Algorithm {
	case "id":
		hash = argon2.IDKey([]byte(password), salt, uint32(aOpts.Iterations), uint32(aOpts.Memory), uint8(aOpts.Parallelism), 32)
	case "i":
		hash = argon2.Key([]byte(password), salt, uint32(aOpts.Iterations), uint32(aOpts.Memory), uint8(aOpts.Parallelism), 32)
	default:
		t.Fatal("not supported algorithm")
	}
	var sharedKey [32]byte
	copy(sharedKey[:], hash)
	nonceRaw, err := keyInfo.Data.SecureNonce.bytes()
	assert.NoError(t, err)
	var nonce [24]byte
	copy(nonce[:], nonceRaw)
	msg, err := keyInfo.Data.SecureBox.bytes()
	assert.NoError(t, err)
	actual, ok := box.OpenAfterPrecomputation([]byte{}, msg, &nonce, &sharedKey)
	assert.True(t, ok)
	assert.Equal(t, priv[:], actual[:])
}
