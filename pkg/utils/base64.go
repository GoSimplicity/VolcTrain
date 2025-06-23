/*
 * Apache License
 * Version 2.0, January 2004
 * http://www.apache.org/licenses/
 *
 * TERMS AND CONDITIONS FOR USE, REPRODUCTION, AND DISTRIBUTION
 *
 * Copyright 2025 Bamboo
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package utils

import (
	"encoding/base64"
	"errors"
)

// Base64Encrypt 用于加密明文密码
func Base64Encrypt(password string) string {
	return base64.StdEncoding.EncodeToString([]byte(password))
}

// Base64Decrypt 用于解密加密后的密码
func Base64Decrypt(encryptedPassword string) (string, error) {
	decoded, err := base64.StdEncoding.DecodeString(encryptedPassword)
	if err != nil {
		return "", errors.New("解密失败: " + err.Error())
	}
	return string(decoded), nil
}

// Base64EncryptWithMagic 加密
// 通过先添加特定盐值、反转字符串再进行base64编码
func Base64EncryptWithMagic(password string) string {
	const salt = "CloudOps@2024#Security!"
	// 添加盐值并反转字符串
	reversed := reverseString(password + salt)
	// 多次编码增加复杂度
	encoded := base64.StdEncoding.EncodeToString([]byte(reversed))
	encoded = base64.StdEncoding.EncodeToString([]byte(encoded + salt))
	return encoded
}

// Base64DecryptWithMagic 解密
// 与加密过程相反的步骤还原原始密码
func Base64DecryptWithMagic(encryptedPassword string) (string, error) {
	const salt = "CloudOps@2024#Security!"
	// 第一次解码
	decoded, err := base64.StdEncoding.DecodeString(encryptedPassword)
	if err != nil {
		return "", errors.New("解密失败: " + err.Error())
	}

	// 移除盐值
	decodedStr := string(decoded)
	if len(decodedStr) <= len(salt) {
		return "", errors.New("解密失败: 无效的加密数据")
	}
	decodedStr = decodedStr[:len(decodedStr)-len(salt)]

	// 第二次解码
	finalDecoded, err := base64.StdEncoding.DecodeString(decodedStr)
	if err != nil {
		return "", errors.New("解密失败: " + err.Error())
	}

	// 反转并移除盐值
	reversed := reverseString(string(finalDecoded))
	if len(reversed) <= len(salt) {
		return "", errors.New("解密失败: 无效的加密数据")
	}

	return reversed[:len(reversed)-len(salt)], nil
}

// 反转字符串
func reverseString(s string) string {
	runes := []rune(s)
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}
	return string(runes)
}
