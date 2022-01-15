package utils

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/mailgun/mailgun-go/v4"
)

func SendLoginCode(recipient string, code string) error {

	// Your available domain names can be found here:
	// (https://app.mailgun.com/app/domains)
	var yourDomain string = os.Getenv("MAILGUN_DOMAIN_NAME")

	// You can find the Private API Key in your Account Menu, under "Settings":
	// (https://app.mailgun.com/app/account/security)
	var privateAPIKey string = os.Getenv("MAILGUN_PRIVATE_API_KEY")

	// Create an instance of the Mailgun Client
	mg := mailgun.NewMailgun(yourDomain, privateAPIKey)

	sender := fmt.Sprintf("mailgun@%s", yourDomain)
	subject := "Your Inception Animals Login Code"
	body := ""

	// The message object allows you to add attachments and Bcc recipients
	message := mg.NewMessage(sender, subject, body, recipient)
	message.SetTemplate("login_code")
	err := message.AddTemplateVariable("code", code)
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	// Send the message with a 10 second timeout
	resp, id, err := mg.Send(ctx, message)

	if err != nil {
		return err
	}

	fmt.Printf("ID: %s Resp: %s\n", id, resp)
	return nil
}
