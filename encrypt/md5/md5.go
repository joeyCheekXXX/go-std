package md5

import (
	"crypto/md5"
	"fmt"
)

func MD5(s string) string {
	sum := md5.Sum([]byte(s))
	return fmt.Sprintf("%x\n", sum)
}
