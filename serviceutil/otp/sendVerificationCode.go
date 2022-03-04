package otp

import (
	"fmt"
	"log"
	"math/rand"
	"time"
)

const otpChars = "1234567890"

func SendVerificationCode(phone string) string {
	code, err := GenerateOTP(6)
	if err != nil {
		log.Fatal(err)
	}
	return code
}

func GenerateEmailVerificationCode() string {
	code, err := GenerateOTP(6)
	if err != nil {
		log.Fatal(err)
	}
	return code
}

func GenerateOTP(length int) (string, error) {
	buffer := make([]byte, length)
	_, err := rand.Read(buffer)
	if err != nil {
		return "", err
	}

	otpCharsLength := len(otpChars)
	for i := 0; i < length; i++ {
		buffer[i] = otpChars[int(buffer[i])%otpCharsLength]
	}

	return string(buffer), nil
}

func GenerateRandomToken() string {
	code := randomString(32)
	return code
}

func randomString(length int) string {
	rand.Seed(time.Now().UnixNano())
	b := make([]byte, length)
	rand.Read(b)
	return fmt.Sprintf("%x", b)[:length]
}
