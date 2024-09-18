package hash_code

import "crypto/sha256"

// hash_code 用于与JAVA项目中的hashCode保持一致

// StringToNumber 将字符串转换为正整数哈希码
func StringToNumber(str string) uint32 {
	// 计算SHA-256哈希值
	digest := sha256.New()
	digest.Write([]byte(str))
	hash := digest.Sum(nil)

	// 将字节转换成正整数
	var hashCode uint32 = 0
	for _, b := range hash {
		hashCode = hashCode*31 + uint32(b&0xff)
		hashCode &= 0x00000000ffffffff // 确保是正数
	}
	return hashCode
}
