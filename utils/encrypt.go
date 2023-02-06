package utils

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/hmac"
	"crypto/sha1"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"hash/crc32"

	"io"
)

func GetHashCode(s string) int64 {
	return int64(crc32.ChecksumIEEE([]byte(s)))
}

func PKCS7Padding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}

func PKCS7UnPadding(origData []byte) []byte {
	length := len(origData)
	unpadding := int(origData[length-1])
	return origData[:(length - unpadding)]
}

//AES加密,CBC
func AesEncrypt(origData, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	blockSize := block.BlockSize()
	origData = PKCS7Padding(origData, blockSize)
	blockMode := cipher.NewCBCEncrypter(block, key[:blockSize])
	crypted := make([]byte, len(origData))
	blockMode.CryptBlocks(crypted, origData)
	return crypted, nil
}

//AES解密
func AesDecrypt(crypted, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	blockSize := block.BlockSize()
	blockMode := cipher.NewCBCDecrypter(block, key[:blockSize])
	origData := make([]byte, len(crypted))
	blockMode.CryptBlocks(origData, crypted)
	origData = PKCS7UnPadding(origData)
	return origData, nil
}

func hashBytes(key string) (hash []byte) {
	h := sha1.New()
	io.WriteString(h, key)
	hashStr := hex.EncodeToString(h.Sum(nil))
	hash = []byte(hashStr)[:32]
	return
}

func CommonAesEncrypt(key, data string) string {
	if len(data) > 0 {
		encrypted, err := AesEncrypt([]byte(data), []byte(key))
		if err == nil {
			return base64.StdEncoding.EncodeToString(encrypted)
		}
	}
	return data
}

func CommonAesDecrypt(key, data string) string {
	if len(data) > 0 {
		encrypted, err := base64.StdEncoding.DecodeString(data)
		if err != nil {
			return data
		}
		tempBytes, err := AesDecrypt(encrypted, []byte(key))
		if err == nil {
			return string(tempBytes)
		}
	}
	return data
}

func EncryptSha256(s string) string {
	hash := sha256.New()
	hash.Write([]byte(s))
	sum := hash.Sum(nil)
	return hex.EncodeToString(sum)
}

func EncryptHmacSha1(secret, text string) (sha string) {
	h := hmac.New(sha1.New, []byte(secret))
	h.Write([]byte(text))
	sha = base64.StdEncoding.EncodeToString(h.Sum(nil))
	return
}