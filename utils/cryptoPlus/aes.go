package cryptoPlus

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"io"
	"prismx_cli/utils/logger"
	"prismx_cli/utils/randomUtils"
)

// PKCS7 填充
func pkcs7Padding(data []byte, blockSize int) []byte {
	padding := blockSize - len(data)%blockSize
	padText := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(data, padText...)
}

// PKCS7 取消填充
func pkcs7Unpadding(data []byte) []byte {
	length := len(data)
	unpadding := int(data[length-1])
	return data[:(length - unpadding)]
}

// AESEncryptCBC CBC 加密
func AESEncryptCBC(plaintext []byte, key []byte) (string, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		logger.Panicf(err.Error())
	}
	plaintext = pkcs7Padding(plaintext, block.BlockSize())
	iv := randomUtils.NewV4().Bytes()
	blockMode := cipher.NewCBCEncrypter(block, iv)
	cipherText := make([]byte, len(plaintext))
	blockMode.CryptBlocks(cipherText, plaintext)
	return base64.StdEncoding.EncodeToString(append(iv[:], cipherText[:]...)), nil
}

// AESDecryptCBC cbc 解密
func AESDecryptCBC(ciphertext, key []byte) ([]byte, error) {
	// 转成字节数组
	crytedByte, _ := base64.StdEncoding.DecodeString(string(ciphertext))
	k := []byte(key)

	// 分组秘钥
	block, _ := aes.NewCipher(k)
	// 获取秘钥块的长度
	blockSize := block.BlockSize()
	// 加密模式
	blockMode := cipher.NewCBCDecrypter(block, k[:blockSize])
	// 创建数组
	orig := make([]byte, len(crytedByte))
	// 解密
	blockMode.CryptBlocks(orig, crytedByte)
	// 去补全码
	orig = pkcs7Unpadding(orig)
	return orig, nil
}

func AESEncryptGCM(Content []byte, key []byte) (string, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}
	nonce := make([]byte, 16)
	_, err = io.ReadFull(rand.Reader, nonce)
	if err != nil {
		return "", err
	}
	aesGCM, err := cipher.NewGCMWithNonceSize(block, 16)
	if err != nil {
		return "", err
	}
	ciphertext := aesGCM.Seal(nil, nonce, Content, nil)
	return base64.StdEncoding.EncodeToString(append(nonce, ciphertext...)), nil
}

func AESDecryptGCM(ciphertext []byte, key []byte) ([]byte, error) {
	cipherBytes, err := base64.StdEncoding.DecodeString(string(ciphertext))
	if err != nil {
		return nil, err
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	nonceSize := 16
	if len(cipherBytes) < nonceSize {
		return nil, errors.New("invalid ciphertext")
	}

	nonce := cipherBytes[:nonceSize]
	encrypted := cipherBytes[nonceSize:]

	aesGCM, err := cipher.NewGCMWithNonceSize(block, nonceSize)
	if err != nil {
		return nil, err
	}

	plaintext, err := aesGCM.Open(nil, nonce, encrypted, nil)
	if err != nil {
		return nil, err
	}
	return plaintext, nil
}
func AESEncryptGCM2(Content []byte, key []byte) (string, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}
	nonce := make([]byte, 16)
	_, err = io.ReadFull(rand.Reader, nonce)
	if err != nil {
		return "", err
	}
	aesGCM, err := cipher.NewGCMWithNonceSize(block, 16)
	if err != nil {
		return "", err
	}
	ciphertext := aesGCM.Seal(nil, nonce, Content, nil)
	return base64.StdEncoding.EncodeToString(append(nonce, ciphertext...)), nil
}
