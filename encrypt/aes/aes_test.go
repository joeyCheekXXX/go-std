package aes

import (
	"fmt"
	"testing"
)

func TestGcmEncrypt(t *testing.T) {

	key := "dskZ/txoPfpmipxF"
	message := `{"id":1, "name":"jobs"}`
	// 加密
	str, _ := Encrypt(key, message)
	// 解密
	str1, _ := Decrypt(key, str)
	// 打印
	fmt.Printf(" 加密：%v\n 解密：%s\n ",
		str, str1,
	)

	// 解密
	str1, _ = Decrypt(key, "DlGBnzwTqbJFjqvttFRI4v+8TXERyQvZEonMpJ6kHSU=")
	// 打印
	fmt.Printf(" 加密：%v\n 解密：%s\n ",
		str, str1,
	)

}
