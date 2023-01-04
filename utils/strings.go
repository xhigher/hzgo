package utils

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"github.com/google/uuid"
)

const (
	allNumbersAndLetters = "0123456789abcdefghijklmnopqrstuvwxyz"
)

func MD5(str string) string {
	h := md5.New()
	h.Write([]byte(str))
	return hex.EncodeToString(h.Sum(nil))
}


func UUID() string {
	return uuid.New().String()
}

func TimeUUID() string {
	return fmt.Sprintf("%s-%s", UUID(), IntToBase36(NowTimeNano()/1000))
}

func IntToBase36(num int64) string {
	result := ""
	size := int64(len(allNumbersAndLetters))
	for num != 0 {
		rem := num % 36
		result = string(allNumbersAndLetters[rem]) + result
		num = num / size
	}
	return result
}