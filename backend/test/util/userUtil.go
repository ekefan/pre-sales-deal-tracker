package test

import (
	"fmt"
	"math/rand"
	"strings"
)

var alphabets = "abcdefghijklmnopqrstuvwxyz"

func RandomString(n int) string {
	var sb strings.Builder
	for i := 1; i <= n; i++ {
		char := alphabets[rand.Intn(len(alphabets))]
		sb.WriteByte(char)
	}
	return sb.String()
}


func GenEmail() string {
	return fmt.Sprintf("%s@gmail.com", RandomString(5))
}

func GenFullname() string {
	return fmt.Sprintf("%s %s", RandomString(5), RandomString(4))
}

func GenPassWord() string {
	return fmt.Sprintf("%s%v", RandomString(6), 10 + rand.Int63n(89))
}
