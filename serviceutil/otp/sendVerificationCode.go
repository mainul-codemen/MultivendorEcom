package otp

import (
	"crypto/rand"
	"log"
)

const otpChars = "1234567890"

func SendVerificationCode(phone string) string {
	code, err := GenerateOTP(6)
	if err != nil {
		log.Fatal(err)
	}
	return code
}
func SendEmailVerificationCode(phone string) string {
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
