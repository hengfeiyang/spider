package util

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"errors"
)

// AesEncrypt aes加密
func AesEncrypt(data, iv, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	blockSize := block.BlockSize()
	data = PKCS5Padding(data, blockSize)
	blockMode := cipher.NewCBCEncrypter(block, iv)
	cryptData := make([]byte, len(data))
	blockMode.CryptBlocks(cryptData, data)
	baseData := make([]byte, base64.StdEncoding.EncodedLen(len(cryptData)))
	base64.StdEncoding.Encode(baseData, cryptData)
	return baseData, nil
}

// AesDecrypt aes解密
func AesDecrypt(data, iv, key []byte) ([]byte, error) {
	baseData := make([]byte, base64.StdEncoding.DecodedLen(len(data)))
	length, err := base64.StdEncoding.Decode(baseData, data)
	if err != nil {
		return nil, err
	}
	data = baseData[:length]
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	blockMode := cipher.NewCBCDecrypter(block, iv)
	blockSize := block.BlockSize()
	origData := make([]byte, len(data))
	blockMode.CryptBlocks(origData, data)
	origData, err = PKCS5UnPadding(origData, blockSize)
	if err != nil {
		return nil, err
	}
	return origData, nil
}

// PKCS5Padding aes加密补码
func PKCS5Padding(data []byte, blockSize int) []byte {
	padding := blockSize - len(data)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(data, padtext...)
}

// PKCS5UnPadding aes解密去码
func PKCS5UnPadding(data []byte, blockSize int) ([]byte, error) {
	length := len(data)
	unpadding := int(data[length-1])
	if unpadding >= blockSize {
		return nil, errors.New("AES PCKS5UnPadding penic, unpadding Illegal")
	}
	return data[:(length - unpadding)], nil
}
