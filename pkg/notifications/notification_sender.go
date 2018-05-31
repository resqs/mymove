package notifications

import (
	"bytes"
	"os"
	"path"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ses"
	"github.com/aws/aws-sdk-go/service/ses/sesiface"
	"github.com/go-gomail/gomail"
	"github.com/pkg/errors"
)

type notification interface {
	emails() ([]emailContent, error)
}

type emailContent struct {
	attachments    []string
	recipientEmail string
	subject        string
	htmlBody       string
	textBody       string
}

// SendNotification sends a one or more notifications for all supported mediums
// TODO: It's starting to look preferrable to refactor configuration of this
// package (SES interface and attachment config).
func SendNotification(notification notification, svc sesiface.SESAPI, attachmentDir string) error {
	emails, err := notification.emails()
	if err != nil {
		return err
	}

	return sendEmails(emails, svc, attachmentDir)
}

func sendEmails(emails []emailContent, svc sesiface.SESAPI, attachmentDir string) error {
	for _, email := range emails {
		rawMessage, err := formatRawEmailMessage(email, attachmentDir)
		if err != nil {
			return err
		}

		input := ses.SendRawEmailInput{
			Destinations: []*string{aws.String(email.recipientEmail)},
			RawMessage:   &ses.RawMessage{Data: rawMessage},
			Source:       aws.String(senderEmail()),
		}

		// Returns the message ID. Should we store that somewhere?
		_, err = svc.SendRawEmail(&input)
		if err != nil {
			return errors.Wrap(err, "Failed to send email using SES")
		}
	}

	return nil
}

func formatRawEmailMessage(email emailContent, attachmentDir string) ([]byte, error) {
	m := gomail.NewMessage()
	m.SetHeader("From", senderEmail())
	m.SetHeader("To", email.recipientEmail)
	m.SetHeader("Subject", email.subject)
	m.SetBody("text/plain", email.textBody)
	m.AddAlternative("text/html", email.htmlBody)
	for _, attachment := range email.attachments {
		m.Attach(path.Join(attachmentDir, attachment))
	}

	buf := new(bytes.Buffer)
	_, err := m.WriteTo(buf)
	if err != nil {
		return buf.Bytes(), errors.Wrap(err, "Failed to generate raw email notification message")
	}

	return buf.Bytes(), nil
}

func senderEmail() string {
	return "noreply@" + os.Getenv("AWS_SES_DOMAIN")
}
