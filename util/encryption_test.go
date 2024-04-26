package util

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestEncryptsAndDecrypts(t *testing.T) {
	a := assert.New(t)

	testSecret := "mysecretpassword"
	testText := "the super secret text"

	encrypted, err := EncryptAndEncode(testText, testSecret)
	a.Nil(err)
	a.NotEmpty(encrypted)
	a.NotEqual(testText, encrypted)

	decrypted, err := DecryptAndDecode(encrypted, testSecret)
	a.Nil(err)
	a.Equal(testText, decrypted)
}

func TestEncodeAndDecode(t *testing.T) {
	a := assert.New(t)

	s := "my to be encoded value"

	encoded, err := ConvertToBase64([]byte(s))
	a.Nil(err)
	a.NotEmpty(encoded)

	decoded, err := ConvertFromBase64(encoded)
	a.Nil(err)
	a.Equal(s, string(decoded))

}
