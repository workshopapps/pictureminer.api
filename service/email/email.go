package utility

import (
	"fmt"
	"github.com/workshopapps/pictureminer.api/utility"
	"net/smtp"
	"os"
)

func EmailSender(senderEmail string, password string, receiverEmail []string, subject string, body string) utility.Response {
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
		rd := utility.BuildErrorResponse(550, "error", "Unable to send email", err, nil)
		return rd
	}
	return utility.BuildSuccessResponse(250, "Email sent successfully", nil)
}
