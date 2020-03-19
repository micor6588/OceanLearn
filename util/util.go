package util

import (
	"math/rand"
	"time"
)

// RandomString 生成10位随机字符串
func RandomString(n int) string {
	var letters = []byte("fhgqewuighweiofwdhiqehqwfbsfiweufowieferpohpkmbncbkxqyyewoevdffovnd")
	result := make([]byte, n)
	rand.Seed(time.Now().Unix())
	for i := range result {
		result[i] = letters[rand.Intn(len(letters))]
	}
	return string(result)
}
