package yzhpay

import (
	"bytes"
	"crypto"
	"crypto/cipher"
	"crypto/des"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"errors"
	"fmt"
)

// Sign 生成签名
func Sign(plaintext, privateKey string) (ciphertext string, err error) {

	priKey, err := LoadPrivateKey([]byte(privateKey))
	if err != nil {
		return
	}

	hash := sha256.New()
	_, err = hash.Write([]byte(plaintext))
	digest := hash.Sum(nil)
	rsaSign, err := rsa.SignPKCS1v15(rand.Reader, priKey, crypto.SHA256, digest)
	if err != nil {
		return
	}
	ciphertext = base64.StdEncoding.EncodeToString(rsaSign)
	return
}

// VerifySign 校验签名
func VerifySign(plaintext, sign, publicKey string) (ok bool, err error) {
	pubKey, err := LoadPublicKey([]byte(publicKey))
	if err != nil {
		return
	}

	b, err := base64.StdEncoding.DecodeString(sign)
	if err != nil {
		return
	}
	hash := sha256.New()
	hash.Write([]byte(plaintext))
	digest := hash.Sum(nil)
	err = rsa.VerifyPKCS1v15(pubKey, crypto.SHA256, digest, b)
	if err != nil {
		return
	}
	ok = true
	return
}

// LoadPrivateKey 加载私钥
func LoadPrivateKey(data []byte) (priv *rsa.PrivateKey, err error) {
	block, _ := pem.Decode(data)
	if block == nil {
		err = fmt.Errorf("decode private key fail")
		return
	}

	privateKey, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		return
	}
	priv, ok := privateKey.(*rsa.PrivateKey)
	if !ok {
		err = errors.New("expected *rsa.PrivateKeytype")
	}
	return
}

// LoadPublicKey 加载公钥
func LoadPublicKey(data []byte) (pub *rsa.PublicKey, err error) {
	block, _ := pem.Decode(data)
	if block == nil {
		err = fmt.Errorf("decode public key fail")
		return
	}

	publicKey, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return
	}
	pub, ok := publicKey.(*rsa.PublicKey)
	if !ok {
		err = errors.New("load public key fail")
	}
	return
}

// Encrypt data加密
func Encrypt(originData []byte, des3key string) (string, error) {
	crypt, err := TripleDesEncrypt(originData, []byte(des3key))
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(crypt), nil
}

// Decrypt data解密
func Decrypt(data string, des3key string) ([]byte, error) {
	crypt, err := base64.StdEncoding.DecodeString(data)
	if err != nil {
		return nil, err
	}
	return TripleDesDecrypt([]byte(crypt), []byte(des3key))
}

// TripleDesEncrypt 3DES加密
func TripleDesEncrypt(originData, des3key []byte) ([]byte, error) {
	block, err := des.NewTripleDESCipher(des3key)
	if err != nil {
		return nil, err
	}
	originData = PKCS5Padding(originData, block.BlockSize())
	blockMode := cipher.NewCBCEncrypter(block, des3key[:8])
	crypt := make([]byte, len(originData))
	blockMode.CryptBlocks(crypt, originData)
	return crypt, nil
}

// TripleDesDecrypt 3DES解密
func TripleDesDecrypt(crypt, des3key []byte) ([]byte, error) {
	block, err := des.NewTripleDESCipher(des3key)
	if err != nil {
		return nil, err
	}
	blockMode := cipher.NewCBCDecrypter(block, des3key[:8])
	originData := make([]byte, len(crypt))
	blockMode.CryptBlocks(originData, crypt)
	originData = PKCS5UnPadding(originData)
	return originData, nil
}

// PKCS5Padding 填充
func PKCS5Padding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}

// PKCS5UnPadding 取消填充
func PKCS5UnPadding(origData []byte) []byte {
	length := len(origData)
	// 去掉最后⼀一个字节 unpadding 次
	unpadding := int(origData[length-1])
	return origData[:(length - unpadding)]
}
