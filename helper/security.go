package helper

import (
	"golang.org/x/crypto/bcrypt"
)

const (
	BcryptCost = 12
	dummyHash  = "$2a$12$e8nZq2f2m4a1g8E2J0c9s.xT0yL7wF2hV4mN6bK8jR1tY3uI5o7Pq"
)

// HashPassword mengenkripsi kata sandi menggunakan bcrypt cost 12
func HashPassword(plain string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(plain), BcryptCost)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

// VerifyPassword membandingkan plain text kata sandi dengan hash bcrypt
func VerifyPassword(hash, plain string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(plain))
	return err == nil
}

// VerifyDummyPassword mengeksekusi komputasi bcrypt identik saat user tidak ditemukan
// untuk mencegah timing attack / user enumeration
func VerifyDummyPassword(plain string) {
	_ = bcrypt.CompareHashAndPassword([]byte(dummyHash), []byte(plain))
}
