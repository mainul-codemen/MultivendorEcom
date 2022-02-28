package mail

import (
	"fmt"
	"net/smtp"

	"github.com/MultivendorEcom/serviceutil/logger"
)

func SendingMail(mail string, otp string) error {

	// Sender data.
	from := "testtune4@gmail.com"
	password := "greenbd1"

	// Receiver email address.
	to := []string{mail}

	// smtp server configuration.
	smtpHost := "smtp.gmail.com"
	smtpPort := "587"

	// Message.
	msg := fmt.Sprintf("This is Email Verification For User. Your OTP is %s", otp)
	message := []byte(msg)

	// Authentication.
	auth := smtp.PlainAuth("", from, password, smtpHost)

	// Sending email.
	err := smtp.SendMail(smtpHost+":"+smtpPort, auth, from, to, message)
	if err != nil {
		msg := fmt.Sprintf("error while sending Email to %s", mail)
		logger.Error(msg)
		return err
	}
	logger.Info("Email Sent Successfully!")
	return nil
}
