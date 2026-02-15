package mailer

import (
	"context"
	"fmt"
	"time"

	"github.com/mailersend/mailersend-go"
	"github.com/neevan0842/BlogSphere/backend/config"
	"go.uber.org/zap"
)

type Mailer struct {
	client    *mailersend.Mailersend
	logger    *zap.SugaredLogger
	fromEmail string
}

func NewMailer(logger *zap.SugaredLogger) *Mailer {
	if config.Envs.MAILERSEND_API_KEY == "" || config.Envs.FROM_EMAIL == "" {
		logger.Warn("MailerSend not configured, mailer will not send emails")
		return &Mailer{
			client:    nil,
			logger:    logger,
			fromEmail: config.Envs.FROM_EMAIL,
		}
	} else {
		return &Mailer{
			client:    mailersend.NewMailersend(config.Envs.MAILERSEND_API_KEY),
			logger:    logger,
			fromEmail: config.Envs.FROM_EMAIL,
		}
	}
}

func (m *Mailer) SendWelcomeEmail(toEmail string, username string) error {
	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	subject := "Welcome to BlogSphere - Start Your Developer Journey!"
	text := fmt.Sprintf(welcomeEmailTextTemplate, username)
	html := fmt.Sprintf(welcomeEmailHTMLTemplate, username)

	from := mailersend.From{
		Name:  "BlogSphere",
		Email: config.Envs.FROM_EMAIL,
	}

	recipients := []mailersend.Recipient{
		{
			Name:  username,
			Email: toEmail,
		},
	}

	message := m.client.Email.NewMessage()

	message.SetFrom(from)
	message.SetRecipients(recipients)
	message.SetSubject(subject)
	message.SetHTML(html)
	message.SetText(text)

	// Send the email with error handling
	res, err := m.client.Email.Send(ctx, message)
	if err != nil {
		m.logger.Errorf("Failed to send welcome email to %s: %v", toEmail, err)
		return err
	}

	m.logger.Infof("Welcome email sent to %s successfully: %v", toEmail, res)
	return nil
}

func (m *Mailer) SendAccountDeletionEmail(toEmail string, username string) error {
	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	subject := "Your BlogSphere Account Has Been Deleted"
	text := fmt.Sprintf(accountDeletionEmailTextTemplate, username)
	html := fmt.Sprintf(accountDeletionEmailHTMLTemplate, username)

	from := mailersend.From{
		Name:  "BlogSphere",
		Email: config.Envs.FROM_EMAIL,
	}

	recipients := []mailersend.Recipient{
		{
			Name:  username,
			Email: toEmail,
		},
	}

	message := m.client.Email.NewMessage()

	message.SetFrom(from)
	message.SetRecipients(recipients)
	message.SetSubject(subject)
	message.SetHTML(html)
	message.SetText(text)

	// Send the email with error handling
	res, err := m.client.Email.Send(ctx, message)
	if err != nil {
		m.logger.Errorf("Failed to send account deletion email to %s: %v", toEmail, err)
		return err
	}

	m.logger.Infof("Account deletion email sent to %s successfully: %v", toEmail, res)
	return nil
}
