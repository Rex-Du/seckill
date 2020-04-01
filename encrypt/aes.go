// Author : rexdu
// Time : 2020-04-01 23:30
package encrypt

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"errors"
)

// 高级加密标准（AES）
// 16,24,32位字符串，分别对应AES-128,AES-192,AES-256三种加密方法
var PwdKey = []byte("DIS**#KKKDJJSKDI") // 不可泄漏

// 三种模式 PKCS7填充模式
func PKCS7Padding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	// 函数的功能是把切片[]byte{byte(padding)}复制padding个，然后合并成新的字节切片返回
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}

func AesEncrypt(origData []byte, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	// 获取块的大小
	blockSize := block.BlockSize()
	// 对数据进行填充，让数据长度满足需求
	origData = PKCS7Padding(origData, blockSize)
	// 采用加密方法中的CBC加密模式
	blockMode := cipher.NewCBCEncrypter(block, key[:blockSize])
	crypted := make([]byte, len(origData))
	blockMode.CryptBlocks(crypted, origData)
	return crypted, nil

}

// 加密base64
func EnPwdCode(pwd []byte) (string, error) {
	result, err := AesEncrypt(pwd, PwdKey)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(result), err
}

// 填充的反向操作，删除填充字符串
func PKCS7UnPadding(origDate []byte) ([]byte, error) {
	length := len(origDate)
	if length == 0 {
		return nil, errors.New("加密字符串错误")
	} else {
		unpadding := int(origDate[length-1])
		if unpadding > length {
			return nil, errors.New("加密字符串错误")
		}
		return origDate[:(length - unpadding)], nil
	}
}

// 实现解密
func AesDeCrypt(cypted []byte, key []byte) ([]byte, error) {
	// 创建加密算法实例
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	// 获取块的大小
	blockSize := block.BlockSize()
	blockMode := cipher.NewCBCDecrypter(block, key[:blockSize])
	origDate := make([]byte, len(cypted))
	// 这个函数也可以用来解密
	blockMode.CryptBlocks(origDate, cypted)
	// 去除填充字符串
	origDate, err = PKCS7UnPadding(origDate)
	if err != nil {
		return nil, err
	}
	return origDate, err
}

// 解密
func DePwdCode(pwd string) ([]byte, error) {
	pwdByte, err := base64.StdEncoding.DecodeString(pwd)
	if err != nil {
		return nil, err
	}
	// 执行AES解密
	return AesDeCrypt(pwdByte, PwdKey)
}
