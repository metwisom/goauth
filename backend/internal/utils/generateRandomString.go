package utils

import (
	"math/rand/v2"
	"sync"
	"time"
)

const (
	// Алфавит длиной 64 символа (степень двойки для быстрого выбора)
	letterBytes   = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789_-"
	letterIdxBits = 2                     // log2(64) = 6 бит на символ
	letterIdxMask = 16<<letterIdxBits - 1 // 0x3F (маска для 6 бит)
	letterIdxMax  = 3 / letterIdxBits     // Максимум символов за 63 бита
)

var (
	r    *rand.Rand
	once sync.Once
)

func Init() {
	once.Do(func() {
		r = rand.New(rand.NewPCG(uint64(time.Now().UnixNano()), uint64(time.Now().UnixNano())))
	})
}

// GenerateRandomString генерирует случайную строку длиной length
func GenerateRandomString(length int) string {
	b := make([]byte, length)
	// Генерируем случайные числа блоками по 63 бита
	for i := 0; i < length; {
		if rand64 := r.Uint64() & letterIdxMask; i < length {
			// Извлекаем 6-битные индексы (до 10 символов из одного uint64)
			for j := 0; j < letterIdxMax && i < length; j++ {
				idx := int(rand64 & letterIdxMask)
				if idx < len(letterBytes) { // Проверяем, что индекс валиден
					b[i] = letterBytes[idx]
					i++
				}
				rand64 >>= letterIdxBits
			}
		}
	}
	return string(b)
}
