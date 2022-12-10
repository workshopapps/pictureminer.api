package utility

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"html/template"
	"log"
	"os"
	"path/filepath"
	"time"

	mail "github.com/xhit/go-simple-mail/v2"
)

type EmailData struct {
	URL       string
	UserName  string
	Subject   string
	Error     string
	BatchName string
}

type MailerConfig struct {
	timeout      time.Duration
	host         string
	port         int
	username     string
	password     string
	sender       string
	templatePath string
}

var tmpl *template.Template

var mailConfig = MailerConfig{
	host:         "email-smtp.us-east-2.amazonaws.com",
	port:         2587,
	templatePath: "templates/*html",
	timeout:      10 * time.Second,
}

func init() {
	tmpl = template.Must(template.ParseGlob(mailConfig.templatePath))
}

func SendMail(from, username, password, receiverEmail string, template string, data *EmailData) error {
	// SMTP Server
	server := mail.SMTPServer{
		Host:           mailConfig.host,
		Port:           mailConfig.port,
		Username:       username,
		Password:       password,
		Encryption:     mail.EncryptionSTARTTLS,
		ConnectTimeout: 10 * time.Second,
		SendTimeout:    10 * time.Second,
		TLSConfig:      &tls.Config{InsecureSkipVerify: true},
	}

	// SMTP client
	smtpClient, err := server.Connect()
	if err != nil {
		log.Fatal(err)
	}

	body, err := parseTemplate(template, data)
	if err != nil {
		return err
	}

	email := mail.NewMSG()
	email.SetFrom(from).AddTo(receiverEmail).SetSubject(data.Subject)
	email.SetBody(mail.TextHTML, body)

	// always check error after send
	if email.Error != nil {
		return email.Error
	}

	// Call Send and pass the client
	err = email.Send(smtpClient)
	if err != nil {
		return err
	}

	return nil
}

func parseTemplate(template string, data *EmailData) (string, error) {
	tmpl, err := tmpl.ParseGlob(mailConfig.templatePath)
	if err != nil {
		return "", err
	}

	body := new(bytes.Buffer)
	err = tmpl.ExecuteTemplate(body, template, data)
	if err != nil {
		return "", err
	}

	return body.String(), nil
}

func ParseTemplateDir(dir string) (*template.Template, error) {
	var paths []string
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			paths = append(paths, path)
		}
		return nil
	})

	fmt.Println("..parsing templates...")

	if err != nil {
		return nil, err
	}

	return template.ParseFiles(paths...)
}
