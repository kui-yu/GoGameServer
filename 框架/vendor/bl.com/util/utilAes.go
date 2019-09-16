package util

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
)

//使用PKCS7进行填充
func pkcs7Padding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}

// 取消PKCS7填充
func pkcs7UnPadding(origData []byte) []byte {
	length := len(origData)
	unpadding := int(origData[length-1])
	return origData[:(length - unpadding)]
}

// AES加密
func aesCBCEncrypt(rawData, aesKey, aesIV []byte) ([]byte, error) {
	// 设置AES密钥
	block, err := aes.NewCipher(aesKey)
	if err != nil {
		return nil, err
	}

	// 填充原文
	blockSize := block.BlockSize()
	rawData = pkcs7Padding(rawData, blockSize)

	// 结果长度
	cipherText := make([]byte, len(rawData))

	//block大小和初始向量大小一定要一致
	mode := cipher.NewCBCEncrypter(block, aesIV)
	mode.CryptBlocks(cipherText, rawData)

	// 返回加密结果
	return cipherText, nil
}

// AES解密
func aesCBCDncrypt(encryptData, aesKey, aesIV []byte) ([]byte, error) {
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
	encryptData = pkcs7UnPadding(encryptData)
	// 返回解密结果
	return encryptData, nil
}

// 返回AES-CBC BASE64 加密结果
func Encrypt(rawData, aesKey, aesIV []byte) (string, error) {
	data, err := aesCBCEncrypt(rawData, aesKey, aesIV)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(data), nil
}

// 返回BASE64 AES-CBC 解密结果
func Dncrypt(rawData string, aesKey []byte, aesIV []byte) (string, error) {
	data, err := base64.StdEncoding.DecodeString(rawData)
	if err != nil {
		return "", err
	}
	dnData, err := aesCBCDncrypt(data, aesKey, aesIV)
	if err != nil {
		return "", err
	}
	return string(dnData), nil
}
