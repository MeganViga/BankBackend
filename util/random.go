package util

import (
	"math/rand"
	"strings"
	"time"
)

func RandomNumber(min, max int)int{
	rand.Seed(int64(time.Now().UnixNano()))
	return rand.Intn(max - min + 1) + min
}


func RandomName(length int)string{
	letters := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	var sb strings.Builder
	for i:= 0; i< length;i++{
		c := letters[RandomNumber(0, len(letters) - 1)]
		sb.WriteByte(c)
	}
	return sb.String()
}

func RandomCurrency()string{
	cur := []string{USD,EUR,CAD}
	return cur[RandomNumber(0, len(cur) - 1)]
}