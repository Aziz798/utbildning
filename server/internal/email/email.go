package email

import (
	"errors"
	"log"
	"os"
	"server/internal/types"

	_ "github.com/joho/godotenv/autoload"
	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

func SendEmail(emailToSend types.Email) error {
	from := mail.NewEmail("Utbildning@miljonbemanning", "Utbildning@miljonbemanning.se")
	to := mail.NewEmail(emailToSend.UserName, emailToSend.EmailTo)
	htmlContent := ""
	message := mail.NewSingleEmail(from, emailToSend.Subject, to, emailToSend.Body, htmlContent)
	client := sendgrid.NewSendClient(os.Getenv("SENDGRID_API_KEY"))
	_, err := client.Send(message)
	if err != nil {
		log.Println(err)
		return errors.New(emailToSend.EmailTo)
	}
	return nil
}
