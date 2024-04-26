package util

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
)

func ConvertToBase64(input []byte) (string, error) {
	if input == nil {
		return "", errors.New("cannot convert to base64 with nil input")
	}
	return base64.URLEncoding.EncodeToString(input), nil
}

func ConvertFromBase64(input string) ([]byte, error) {
	if input == "" {
		return nil, errors.New("cannot convert from base64 with blank input")
	}
	return base64.URLEncoding.DecodeString(input)
}

// EncryptAndEncode encrypts an input with a secret key and returns the base64 encoded and encrypted string
func EncryptAndEncode(input string, secretKey string) (string, error) {
	var err error
	var block cipher.Block
	if block, err = aes.NewCipher([]byte(secretKey)); err != nil {
		return "", fmt.Errorf("cannot create cipher for encrypting: %w", err)
	}

	var gcm cipher.AEAD
	if gcm, err = cipher.NewGCM(block); err != nil {
		return "", fmt.Errorf("cannot create gcm for encrypting: %w", err)
	}

	nonce := make([]byte, gcm.NonceSize())

	if _, err = rand.Read(nonce); err != nil {
		return "", fmt.Errorf("cannot read nonce for encrypting: %w", err)
	}

	ciphertext := gcm.Seal(nonce, nonce, []byte(input), nil)

	var encoded string
	if encoded, err = ConvertToBase64(ciphertext); err != nil {
		return "", fmt.Errorf("cannot encode encrypted input: %w", err)
	}

	return encoded, nil
}

// DecryptAndDecode decodes and decrypts a given base64 string and returns the decrypted plain text
func DecryptAndDecode(base64Encoded string, secretKey string) (string, error) {
	var decoded []byte
	var err error
	if decoded, err = ConvertFromBase64(base64Encoded); err != nil {
		return "", fmt.Errorf("cannot decode encoded input: %w", err)
	}

	decodedCipherText := string(decoded)

	var block cipher.Block
	if block, err = aes.NewCipher([]byte(secretKey)); err != nil {
		return "", fmt.Errorf("cannot create cipher for decrypting: %w", err)
	}

	var gcm cipher.AEAD
	if gcm, err = cipher.NewGCM(block); err != nil {
		return "", fmt.Errorf("cannot create gcm for decrypting: %w", err)
	}

	nonceSize := gcm.NonceSize()
	nonce, decodedCipherText := decodedCipherText[:nonceSize], decodedCipherText[nonceSize:]

	var plaintext []byte
	if plaintext, err = gcm.Open(nil, []byte(nonce), []byte(decodedCipherText), nil); err != nil {
		return "", fmt.Errorf("cannot decrypt: %w", err)
	}

	return string(plaintext), nil
}
