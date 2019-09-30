package main

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
)

var aesKey = []byte("0123456789abcdef")
var aesIV = []byte("abcdef0123456789")

//使用PKCS7进行填充，IOS也是7
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

//aes加密，填充秘钥key的16位，24,32分别对应AES-128, AES-192, or AES-256.
func AesCBCEncrypt(rawData []byte) ([]byte, error) {
	block, err := aes.NewCipher(aesKey)
	if err != nil {
		panic(err)
	}

	// 填充原文
	blockSize := block.BlockSize()
	rawData = PKCS7Padding(rawData, blockSize)

	// 结果长度
	cipherText := make([]byte, len(rawData))

	//block大小和初始向量大小一定要一致
	mode := cipher.NewCBCEncrypter(block, aesIV)
	mode.CryptBlocks(cipherText, rawData)

	return cipherText, nil
}

func AesCBCDncrypt(encryptData []byte) ([]byte, error) {
	block, err := aes.NewCipher(aesKey)
	if err != nil {
		panic(err)
	}

	blockSize := block.BlockSize()

	if len(encryptData) < blockSize {
		panic("ciphertext too short")
	}

	// CBC mode always works in whole blocks.
	if len(encryptData)%blockSize != 0 {
		panic("ciphertext is not a multiple of the block size")
	}

	mode := cipher.NewCBCDecrypter(block, aesIV)

	// CryptBlocks can work in-place if the two arguments are the same.
	mode.CryptBlocks(encryptData, encryptData)
	//解填充
	encryptData = PKCS7UnPadding(encryptData)
	return encryptData, nil
}

func Encrypt(rawData []byte) (string, error) {
	data, err := AesCBCEncrypt(rawData)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(data), nil
}

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
