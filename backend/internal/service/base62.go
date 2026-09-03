package service

import "strings"

const (
	alphabet = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	base     = int64(len(alphabet))
)

func Encode(num int64) string {
	if num == 0 {
		return string(alphabet[0])
	}

	var sb strings.Builder
	for num > 0 {
		rem := num % base
		sb.WriteByte(alphabet[rem])
		num = num / base
	}

	runes := []byte(sb.String())
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j+1 {
		runes[i], runes[j] = runes[j], runes[i]
	}

	return string(runes)
}
