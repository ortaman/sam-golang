package repository

import (
	"bytes"
	"fmt"
	"net/smtp"
	"os"
	"text/template"

	"github.com/ortaman/stori-test/entity"
)

func SendEmail(email_to []string, termplateDir string, transactionResume *entity.TransactionResume) error {
	from := os.Getenv("EMAIL_FROM")
	password := os.Getenv("EMAIL_CODE")

	smtpHost := os.Getenv("SMTP_HOST")
	smtpPort := os.Getenv("SMTP_PORT")

	// Message
	var body bytes.Buffer

	mimeHeaders := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"
	body.Write([]byte(fmt.Sprintf("Subject: Stori Test \n%s\n\n", mimeHeaders)))

	template, TemplatError := template.ParseFiles(termplateDir)

	if TemplatError != nil {
		return fmt.Errorf("%s", TemplatError.Error())
	}

	template.Execute(&body, transactionResume)

	// Authentication and send email
	auth := smtp.PlainAuth("", from, password, smtpHost)
	emailError := smtp.SendMail(smtpHost+":"+smtpPort, auth, from, email_to, body.Bytes())

	if emailError != nil {
		return fmt.Errorf("%s", emailError.Error())
	}

	return nil
}
