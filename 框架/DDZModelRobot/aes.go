package main

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
)

// 填充秘钥key的16位 AES128
var aesKey = []byte("0123456789abcdef")

// 填充向量iv的16位
var aesIV = []byte("abcdef0123456789")

//使用PKCS7进行填充
func PKCS7Padding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}

// 取消PKCS7填充
func PKCS7UnPadding(origData []byte) []byte {
	length := len(origData)
	unpadding := int(origData[length-1])
	return origData[:(length - unpadding)]
}

// AES加密
func AesCBCEncrypt(rawData []byte) ([]byte, error) {
	// 设置AES密钥
	block, err := aes.NewCipher(aesKey)
	if err != nil {
		return nil, err
	}

	// 填充原文
	blockSize := block.BlockSize()
	rawData = PKCS7Padding(rawData, blockSize)

	// 结果长度
	cipherText := make([]byte, len(rawData))

	//block大小和初始向量大小一定要一致
	mode := cipher.NewCBCEncrypter(block, aesIV)
	mode.CryptBlocks(cipherText, rawData)

	// 返回加密结果
	return cipherText, nil
}

// AES解密
func AesCBCDncrypt(encryptData []byte) ([]byte, error) {
	// 设置AES密钥
	block, err := aes.NewCipher(aesKey)
	if err != nil {
		return nil, err
	}

	blockSize := block.BlockSize()

	if len(encryptData) < blockSize {
		return nil, nil
	}

	// AES-CBC加密结果长度为 blockSize 的整数倍
	if len(encryptData)%blockSize != 0 {
		return nil, nil
	}

	mode := cipher.NewCBCDecrypter(block, aesIV)

	// 解密，允许和原文使用同一个变量，会自动覆盖
	mode.CryptBlocks(encryptData, encryptData)

	// 解填充
	encryptData = PKCS7UnPadding(encryptData)
	// 返回解密结果
	return encryptData, nil
}

// 返回AES-CBC BASE64 加密结果
func Encrypt(rawData []byte) (string, error) {
	data, err := AesCBCEncrypt(rawData)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(data), nil
}

// 返回BASE64 AES-CBC 解密结果
func Dncrypt(rawData string) (string, error) {
	data, err := base64.StdEncoding.DecodeString(rawData)
	if err != nil {
		return "", err
	}
	dnData, err := AesCBCDncrypt(data)
	if err != nil {
		return "", err
	}
	return string(dnData), nil
}
