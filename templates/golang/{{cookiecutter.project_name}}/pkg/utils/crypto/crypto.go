// Package crypto 提供各类算法加解密
package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/hex"
	"io"

	"github.com/pkg/errors"
)

// AESEncrypt 实现 AES-GCM 算法加密
func AESEncrypt(key, data string) (string, error) {
	block, err := aes.NewCipher([]byte(key))
	if err != nil {
		return "", errors.Wrap(err, "NewCipher")
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", errors.Wrap(err, "NewGCM")
	}
	// 生成随机因子
	nonce := make([]byte, gcm.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return "", errors.Wrap(err, "make random nonce")
	}
	seal := gcm.Seal(nonce, nonce, []byte(data), nil)
	return hex.EncodeToString(seal), nil
}

// AESDecrypt 实现 AES-GCM 算法解密
func AESDecrypt(key, data string) (string, error) {
	dataByte, err := hex.DecodeString(data)
	if err != nil {
		return "", errors.Wrap(err, "hex decode string")
	}
	block, err := aes.NewCipher([]byte(key))
	if err != nil {
		return "", errors.Wrap(err, "NewCipher")
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", errors.Wrap(err, "NewGCM")
	}
	nonceSize := gcm.NonceSize()
	if len(dataByte) < nonceSize {
		return "", errors.Errorf("ciphertext too short, at least %d", nonceSize)
	}

	nonce, ciphertext := dataByte[:nonceSize], dataByte[nonceSize:]
	open, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return "", errors.Wrap(err, "gcm open")
	}
	return string(open), nil
}
