package main

import (
	"bytes"
	"errors"
	"text/template"

	mail "github.com/xhit/go-simple-mail/v2"
)

var (
	successMessageTemplate = "Dear {{.}},\n Thank you for registering with us. We are glad to have you on board. \n Regards, \n HNG Team"
	t                      = createTemplate("t2", successMessageTemplate)
)

func sendEmail(email string, message string) error {
	server := mail.NewSMTPClient()
	server.Host = "smtp.mailtrap.io"   // change this to take variable from config
	server.Port = 2525                 // change this to take variable from config
	server.Username = "f59d2513ff6024" // change this to take variable from config
	server.Password = "8f448ea24f5f0e" // change this to take variable from config
	server.Encryption = mail.EncryptionTLS

	smtpClient, err := server.Connect()
	if err != nil {
		return err
	}

	// Create email
	emsg := mail.NewMSG()
	emsg.SetFrom("from@example.com") // change this to take variable from config
	emsg.AddTo(email)
	emsg.SetSubject("Sign up registration successful")

	emsg.SetBody(mail.TextHTML, message)

	// Send email
	err = emsg.Send(smtpClient)
	if err != nil {
		return err
	}

	return nil
}
func createTemplate(name, t string) *template.Template {
	return template.Must(template.New(name).Parse(t))
}

func SendSuccessMessage(name string, email string) error {
	// validate Inputs
	if email == "" {
		return errors.New("email is empty")
	}
	if name == "" {
		return errors.New("name is empty")
	}

	// generate message
	var successMessage bytes.Buffer
	t.Execute(&successMessage, name)

	// send email
	err := sendEmail(email, successMessage.String())

	// error handling
	if err != nil {
		return err
	}
	return nil
}
