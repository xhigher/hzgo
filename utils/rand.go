package utils

import "math/rand"

func init() {
	rand.Seed(NowTimeNano())
}

func RandInt64(start, end int64) int64 {
	return rand.Int63n(end-start) + start
}

func RandInt32(start, end int32) int32 {
	return rand.Int31n(end-start) + start
}

func RandString(size int) string {
	bytes := make([]byte, size)
	for i := 0; i < size; i++ {
		b := rand.Intn(10) + 48
		bytes[i] = byte(b)
	}
	return string(bytes)
}

func RandNumberString(size int) string {
	bytes := make([]byte, size)
	for i := 0; i < size; i++ {
		b := rand.Intn(10) + 48
		bytes[i] = byte(b)
	}
	return string(bytes)
}

func RandLetterString(size int) string {
	bytes := make([]byte, size)
	for i := 0; i < size; i++ {
		b := rand.Intn(26) + 65
		bytes[i] = byte(b)
	}
	return string(bytes)
}

