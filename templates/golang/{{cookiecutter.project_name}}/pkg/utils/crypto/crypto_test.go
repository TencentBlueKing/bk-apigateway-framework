package crypto_test

import (
	"crypto/rand"
	"io"
	"testing"

	"github.com/stretchr/testify/assert"

	"bk.tencent.com/{{cookiecutter.project_name}}/{{cookiecutter.project_name}}/pkg/utils/crypto"
)

func TestEncryptDecrypt(t *testing.T) {
	tests := []struct {
		name    string
		encFunc func(key, data string) (string, error)
		decFunc func(key, data string) (string, error)
		keySize int
	}{
		{
			name:    "AES",
			encFunc: crypto.AESEncrypt,
			decFunc: crypto.AESDecrypt,
			keySize: 32,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			plaintext := "hello world"
			key := make([]byte, tt.keySize)
			if _, err := io.ReadFull(rand.Reader, key); err != nil {
				t.Fatal(err)
			}
			// 第一次加密
			ciphertext1, err := tt.encFunc(string(key), plaintext)
			assert.Nil(t, err)

			// 第二次加密
			ciphertext2, err := tt.encFunc(string(key), plaintext)
			assert.Nil(t, err)

			// 确保两次加密后的结果不一致
			assert.NotEqual(t, ciphertext1, ciphertext2)

			// 解密第一次加密的结果
			decrypted1, err := tt.decFunc(string(key), ciphertext1)
			assert.Nil(t, err)
			assert.Equal(t, plaintext, decrypted1)

			// 解密第二次加密的结果
			decrypted2, err := tt.decFunc(string(key), ciphertext2)
			assert.Nil(t, err)
			assert.Equal(t, plaintext, decrypted2)
		})
	}
}
