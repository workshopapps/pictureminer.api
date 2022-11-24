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
	server.Host = "smtp.mailtrap.io"   
	server.Port = 2525                 
	server.Username = "f59d2513ff6024" 
	server.Password = "8f448ea24f5f0e" 
	server.Encryption = mail.EncryptionTLS

	smtpClient, err := server.Connect()
	if err != nil {
		return err
	}

	
	emsg := mail.NewMSG()
	emsg.SetFrom("from@example.com") 
	emsg.AddTo(email)
	emsg.SetSubject("Sign up registration successful")

	emsg.SetBody(mail.TextHTML, message)

	
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
	
	if email == "" {
		return errors.New("email is empty")
	}
	if name == "" {
		return errors.New("name is empty")
	}

	
	var successMessage bytes.Buffer
	t.Execute(&successMessage, name)

	
	err := sendEmail(email, successMessage.String())

	
	if err != nil {
		return err
	}
	return nil
}
