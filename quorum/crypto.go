package quorum

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"

	"golang.org/x/crypto/argon2"
	"golang.org/x/crypto/nacl/box"
)

const (
	saltLength = 16
	nonceLenth = 24
)

var (
	defaultArgonOpts = argonOptions{
		Algorithm:   "i",
		Iterations:  10,
		Memory:      1048576,
		Parallelism: 4,
	}
)

type argonOptions struct {
	Algorithm   string `json:"variant"`
	Iterations  int    `json:"iterations"`
	Memory      int    `json:"memory"`
	Parallelism int    `json:"parallelism"`
}

type StandardBase64EncodedString string

func toStandardBase64EncodedString(src []byte) StandardBase64EncodedString {
	return StandardBase64EncodedString(base64.StdEncoding.EncodeToString(src))
}

func (s StandardBase64EncodedString) bytes() ([]byte, error) {
	return base64.StdEncoding.DecodeString(string(s))
}

type keyData struct {
	Config *privateKeyInfo             `json:"config"`
	PubKey StandardBase64EncodedString `json:"publicKey"`
}

type privateKeyInfo struct {
	Data privateKeyData `json:"data"`
	Type string         `json:"type"`
}

type privateKeyData struct {
	Value       StandardBase64EncodedString `json:"bytes,omitempty"`
	SecureNonce StandardBase64EncodedString `json:"snonce,omitempty"`
	ArgonSalt   StandardBase64EncodedString `json:"asalt,omitempty"`
	SecureBox   StandardBase64EncodedString `json:"sbox,omitempty"`
	ArgonOpts   *argonOptions               `json:"aopts,omitempty"`
}

// first encrypt the privateKey with random salt, random nonce and the password
// then return the JSON reprentation of the key data
func toKeyDataJSON(password string, aOpts *argonOptions, privateKey []byte, pubKeyB64 StandardBase64EncodedString) (string, string, error) {
	salt := make([]byte, saltLength)
	if _, err := rand.Read(salt); err != nil {
		return "", "", fmt.Errorf("random salt error: %s", err)
	}
	var nonce [nonceLenth]byte
	if _, err := rand.Read(nonce[:]); err != nil {
		return "", "", fmt.Errorf("random nonce error: %s", err)
	}
	var info *privateKeyInfo
	if password == "" {
		info = &privateKeyInfo{
			Data: privateKeyData{
				Value: toStandardBase64EncodedString(privateKey),
			},
			Type: "unlocked",
		}
	} else {
		if aOpts == nil {
			aOpts = &defaultArgonOpts
		}
		var hash []byte
		switch aOpts.Algorithm {
		case "id":
			hash = argon2.IDKey([]byte(password), salt, uint32(aOpts.Iterations), uint32(aOpts.Memory), uint8(aOpts.Parallelism), 32)
		case "i":
			hash = argon2.Key([]byte(password), salt, uint32(aOpts.Iterations), uint32(aOpts.Memory), uint8(aOpts.Parallelism), 32)
		default:
			return "", "", fmt.Errorf("unsupported algorithm")
		}
		var sharedKey [32]byte
		copy(sharedKey[:], hash)
		encryptedPrivateKey := box.SealAfterPrecomputation([]byte{}, privateKey, &nonce, &sharedKey)
		info = &privateKeyInfo{
			Data: privateKeyData{
				SecureNonce: toStandardBase64EncodedString(nonce[:]),
				ArgonSalt:   toStandardBase64EncodedString(salt),
				SecureBox:   toStandardBase64EncodedString(encryptedPrivateKey),
				ArgonOpts:   aOpts,
			},
			Type: "argon2sbox",
		}
	}
	rawJson, err := json.Marshal(&keyData{
		Config: info,
		PubKey: pubKeyB64,
	})
	if err != nil {
		return "", "", err
	}
	privateKeyJson, err := json.Marshal(info)
	if err != nil {
		return "", "", err
	}
	return string(rawJson), string(privateKeyJson), nil
}
