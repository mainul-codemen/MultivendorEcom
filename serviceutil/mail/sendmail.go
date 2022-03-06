package mail

import (
	"bytes"
	"fmt"
	"html/template"
	"net/smtp"
	"os"
)

type MailStruct struct {
	From               string
	To                 []string
	Message            string
	Token              string
	OTP                string
	ResetPasswordLinks string
	Subject            string
	UserID             string
}

var tn = "mailtemp.html"

func SendingMail(mail string, ms MailStruct) error {
	from := "testtune4@gmail.com"
	password := "greenbd1"
	to := []string{mail}
	smtpHost := "smtp.gmail.com"
	smtpPort := "587"
	auth := smtp.PlainAuth("", from, password, smtpHost)
	gwd, _ := os.Getwd()
	dir := fmt.Sprintf("%s/serviceutil/mail/%s", gwd, tn)
	t, _ := template.ParseFiles(dir)
	var body bytes.Buffer
	mimeHeaders := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"
	body.Write([]byte(fmt.Sprintf("Subject:%s\n%s\n\n", ms.Subject, mimeHeaders)))
	t.Execute(&body, MailStruct{
		Message:            ms.Message,
		Token:              tn,
		OTP:                ms.OTP,
		ResetPasswordLinks: ms.ResetPasswordLinks,
		UserID:             ms.UserID,
	})

	err := smtp.SendMail(smtpHost+":"+smtpPort, auth, from, to, body.Bytes())
	if err != nil {
		return err

	}
	fmt.Println("Email Sent!")
	return nil
}
