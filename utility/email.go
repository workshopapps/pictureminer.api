package utility

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"html/template"
	"log"
	"net/smtp"
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

func SendMailNaked(from, username, password, receiverEmail string, template string, data *EmailData) error {
	host, port := "smtp.gmail.com", "465"
	serverName := host + ":" + port

	body, err := parseTemplate(template, data)
	if err != nil {
		fmt.Println(err)
		return err
	}

	auth := smtp.PlainAuth("", from, password, host)

	tlsconfig := &tls.Config{
		InsecureSkipVerify: false,
		ServerName:         host,
	}
	conn, err := tls.Dial("tcp", serverName, tlsconfig)
	if err != nil {
		fmt.Println(2, err)
		return err
	}

	client, err := smtp.NewClient(conn, host)
	if err != nil {
		fmt.Println(3, err)
	}
	client.Auth(auth)

	client.Mail(from)
	client.Rcpt(receiverEmail)

	w, err := client.Data()
	if err != nil {
		return err
	}

	w.Write([]byte("subject:" + data.Subject + "\n"))
	w.Write([]byte("Content-type: text/html; charset=\"UTF-8\"\n"))
	w.Write([]byte(body))
	w.Close()

	client.Quit()

	return nil
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
	fmt.Println("here", err)
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
