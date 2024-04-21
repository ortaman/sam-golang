package repository

import (
	"bytes"
	"fmt"
	"net/smtp"
	"text/template"

	"github.com/ortaman/stori-test/entity"
)

type EmailRepository struct {
	PSQLConfig *map[string]string
}

func NewEmailRepository(PSQLConfig *map[string]string) entity.EmailRepoI {
	return &EmailRepository{PSQLConfig}
}

func (emailRepository *EmailRepository) SendEmail(transactionsResume *entity.TransactionResume, emails_to []string, termplateDir string) error {
	// Message
	var body bytes.Buffer

	mimeHeaders := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"
	body.Write([]byte(fmt.Sprintf("Subject: Stori Test \n%s\n\n", mimeHeaders)))

	template, TemplatError := template.ParseFiles(termplateDir)

	if TemplatError != nil {
		return fmt.Errorf("unable to parse email template. %s", TemplatError.Error())
	}

	if err := template.Execute(&body, *transactionsResume); err != nil {
		return fmt.Errorf("unable to execute email html. %v", err)
	}

	Emailconfig := *emailRepository.PSQLConfig

	// Authentication and send email
	auth := smtp.PlainAuth("", Emailconfig["from"], Emailconfig["password"], Emailconfig["smtpHost"])

	emailError := smtp.SendMail(
		Emailconfig["smtpHost"]+":"+Emailconfig["smtpPort"],
		auth,
		Emailconfig["from"],
		emails_to,
		body.Bytes())

	if emailError != nil {
		return fmt.Errorf("send email failed. %s", emailError.Error())
	}

	return nil
}
