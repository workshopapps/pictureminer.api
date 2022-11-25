package utility

import (
	"bytes"
	"fmt"
	"html/template"
	"log"
	"net/smtp"
	"os"
	"path/filepath"

	"github.com/workshopapps/pictureminer.api/internal/model"
)

func EmailSender(senderEmail string, password string, receiverEmail []string, subject string, body string) Response {
	//connect
	host := os.Getenv("host")
	port := os.Getenv("port")
	address := host + ":" + port

	//	message
	message := []byte(fmt.Sprintf("To: %s \r\n"+"Subject: %s \r\n"+"%s", receiverEmail, subject, body))
	//Authentication
	auth := smtp.PlainAuth("", senderEmail, password, host)

	//Send Email
	err := smtp.SendMail(address, auth, senderEmail, receiverEmail, message)
	if err != nil {
		rd := BuildErrorResponse(550, "error", "Unable to send email", err, nil)
		return rd
	}
	return BuildSuccessResponse(250, "Email sent successfully", nil)
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


type EmailData struct {
	URL       string
	FirstName string
	Subject   string
}

// ? Email template parser


func SendEmail(user *model.User, data *EmailData, sender, password string, templateName string) error {
	//EmailSender(Sender, user.Password, )

	var body bytes.Buffer

	template, err := ParseTemplateDir("templates")
	if err != nil {
		log.Fatal("Could not parse template", err)
	}

	template = template.Lookup(templateName)
	template.Execute(&body, &data)
	fmt.Println(template.Name())

	EmailSender(sender, )
}
