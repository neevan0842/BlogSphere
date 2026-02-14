package utils

import (
	"context"
	"fmt"
	"time"

	"github.com/mailersend/mailersend-go"
	"github.com/neevan0842/BlogSphere/backend/config"
	"go.uber.org/zap"
)

func SendWelcomeEmail(toEmail string, username string, logger *zap.SugaredLogger) error {
	if config.Envs.MAILERSEND_API_KEY == "" || config.Envs.FROM_EMAIL == "" {
		logger.Warn("MailerSend not configured, skipping welcome email")
		return nil
	}

	// Create an instance of the mailersend client
	ms := mailersend.NewMailersend(config.Envs.MAILERSEND_API_KEY)

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

	message := ms.Email.NewMessage()

	message.SetFrom(from)
	message.SetRecipients(recipients)
	message.SetSubject(subject)
	message.SetHTML(html)
	message.SetText(text)

	// Send the email with error handling
	res, err := ms.Email.Send(ctx, message)
	if err != nil {
		logger.Errorf("Failed to send welcome email to %s: %v", toEmail, err)
		return err
	}

	logger.Infof("Welcome email sent to %s successfully: %v", toEmail, res)
	return nil
}

const welcomeEmailTextTemplate = `Hi %s,

Welcome to BlogSphere! We're thrilled to have you join our community of developers and tech enthusiasts.

BlogSphere is your platform to share knowledge, connect with fellow developers, and grow your presence in the tech community. Whether you're here to write, read, or engage with content, we're excited to see what you'll create.

Here's what you can do next:
- Complete your profile to personalize your presence
- Explore trending articles from developers worldwide
- Start writing your first blog post and share your expertise
- Engage with the community through comments and likes

If you have any questions or need assistance, feel free to reach out. We're here to help!

Happy blogging,
The BlogSphere Team`

const welcomeEmailHTMLTemplate = `<!DOCTYPE html>
<html>
<head>
	<meta charset="UTF-8">
	<meta name="viewport" content="width=device-width, initial-scale=1.0">
	<style>
		body { font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, 'Helvetica Neue', Arial, sans-serif; line-height: 1.6; color: #333; margin: 0; padding: 0; background-color: #f4f4f5; }
		.container { max-width: 600px; margin: 0 auto; background-color: #ffffff; border-radius: 8px; overflow: hidden; box-shadow: 0 2px 8px rgba(0,0,0,0.1); }
		.header { background: linear-gradient(135deg, #667eea 0%%, #764ba2 100%%); padding: 40px 20px; text-align: center; }
		.header h1 { color: #ffffff; margin: 0; font-size: 28px; font-weight: 700; }
		.content { padding: 40px 30px; }
		.greeting { font-size: 18px; color: #18181b; margin-bottom: 20px; }
		.message { color: #52525b; font-size: 16px; line-height: 1.8; margin-bottom: 30px; }
		.features { background-color: #f9fafb; border-radius: 8px; padding: 20px; margin: 30px 0; }
		.features h3 { color: #18181b; font-size: 16px; margin: 0 0 15px 0; font-weight: 600; }
		.features ul { margin: 0; padding-left: 20px; color: #52525b; }
		.features li { margin-bottom: 10px; }
		.cta { text-align: center; margin: 30px 0; }
		.button { display: inline-block; padding: 14px 32px; background: linear-gradient(135deg, #667eea 0%%, #764ba2 100%%); color: #ffffff; text-decoration: none; border-radius: 6px; font-weight: 600; font-size: 16px; }
		.footer { background-color: #f9fafb; padding: 30px; text-align: center; border-top: 1px solid #e4e4e7; }
		.footer p { color: #71717a; font-size: 14px; margin: 5px 0; }
		.footer a { color: #667eea; text-decoration: none; }
	</style>
</head>
<body>
	<div class="container">
		<div class="header">
			<h1>🚀 Welcome to BlogSphere</h1>
		</div>
		<div class="content">
			<p class="greeting">Hi <strong>%s</strong>,</p>
			<p class="message">
				Welcome to <strong>BlogSphere</strong>! We're thrilled to have you join our community of developers and tech enthusiasts.
			</p>
			<p class="message">
				BlogSphere is your platform to share knowledge, connect with fellow developers, and grow your presence in the tech community. Whether you're here to write, read, or engage with content, we're excited to see what you'll create.
			</p>
			<div class="features">
				<h3>🎯 Here's what you can do next:</h3>
				<ul>
					<li>Complete your profile to personalize your presence</li>
					<li>Explore trending articles from developers worldwide</li>
					<li>Start writing your first blog post and share your expertise</li>
					<li>Engage with the community through comments and likes</li>
				</ul>
			</div>
			<p class="message">
				If you have any questions or need assistance, feel free to reach out. We're here to help!
			</p>
			<p class="message" style="margin-top: 30px; color: #18181b; font-weight: 500;">
				Happy blogging,<br>
				The BlogSphere Team
			</p>
		</div>
		<div class="footer">
			<p>This email was sent to you because you created an account on BlogSphere.</p>
			<p>© 2026 BlogSphere. All rights reserved.</p>
		</div>
	</div>
</body>
</html>`

func SendAccountDeletionEmail(toEmail string, username string, logger *zap.SugaredLogger) error {
	if config.Envs.MAILERSEND_API_KEY == "" || config.Envs.FROM_EMAIL == "" {
		logger.Warn("MailerSend not configured, skipping account deletion email")
		return nil
	}

	// Create an instance of the mailersend client
	ms := mailersend.NewMailersend(config.Envs.MAILERSEND_API_KEY)
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

	message := ms.Email.NewMessage()

	message.SetFrom(from)
	message.SetRecipients(recipients)
	message.SetSubject(subject)
	message.SetHTML(html)
	message.SetText(text)

	// Send the email with error handling
	res, err := ms.Email.Send(ctx, message)
	if err != nil {
		logger.Errorf("Failed to send account deletion email to %s: %v", toEmail, err)
		return err
	}

	logger.Infof("Account deletion email sent to %s successfully: %v", toEmail, res)
	return nil
}

const accountDeletionEmailTextTemplate = `Hi %s,

This email confirms that your BlogSphere account has been successfully deleted from our platform.

All your data, including:
- Your profile information
- Published blog posts
- Comments and interactions
- Account preferences

has been permanently removed from our servers. We're sorry to see you go, but we respect your decision.

If this deletion was made in error or you have any concerns, please contact our support team immediately. Note that account recovery may not be possible after deletion.

We'd love to hear your feedback about your experience on BlogSphere. Your insights help us improve our platform for the developer community.

Thank you for being part of BlogSphere. We hope to see you again in the future!

Best regards,
The BlogSphere Team`

const accountDeletionEmailHTMLTemplate = `<!DOCTYPE html>
<html>
<head>
	<meta charset="UTF-8">
	<meta name="viewport" content="width=device-width, initial-scale=1.0">
	<style>
		body { font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, 'Helvetica Neue', Arial, sans-serif; line-height: 1.6; color: #333; margin: 0; padding: 0; background-color: #f4f4f5; }
		.container { max-width: 600px; margin: 0 auto; background-color: #ffffff; border-radius: 8px; overflow: hidden; box-shadow: 0 2px 8px rgba(0,0,0,0.1); }
		.header { background: linear-gradient(135deg, #ef4444 0%%, #dc2626 100%%); padding: 40px 20px; text-align: center; }
		.header h1 { color: #ffffff; margin: 0; font-size: 26px; font-weight: 700; }
		.content { padding: 40px 30px; }
		.greeting { font-size: 18px; color: #18181b; margin-bottom: 20px; }
		.message { color: #52525b; font-size: 16px; line-height: 1.8; margin-bottom: 20px; }
		.info-box { background-color: #fef2f2; border-left: 4px solid #ef4444; border-radius: 4px; padding: 20px; margin: 25px 0; }
		.info-box h3 { color: #dc2626; font-size: 16px; margin: 0 0 15px 0; font-weight: 600; }
		.info-box ul { margin: 0; padding-left: 20px; color: #52525b; }
		.info-box li { margin-bottom: 8px; }
		.warning { background-color: #fffbeb; border-left: 4px solid #f59e0b; border-radius: 4px; padding: 15px; margin: 20px 0; color: #92400e; font-size: 14px; }
		.footer { background-color: #f9fafb; padding: 30px; text-align: center; border-top: 1px solid #e4e4e7; }
		.footer p { color: #71717a; font-size: 14px; margin: 5px 0; }
		.footer a { color: #667eea; text-decoration: none; }
	</style>
</head>
<body>
	<div class="container">
		<div class="header">
			<h1>Account Deletion Confirmation</h1>
		</div>
		<div class="content">
			<p class="greeting">Hi <strong>%s</strong>,</p>
			<p class="message">
				This email confirms that your BlogSphere account has been <strong>successfully deleted</strong> from our platform.
			</p>
			<div class="info-box">
				<h3>🗑️ The following data has been permanently removed:</h3>
				<ul>
					<li>Your profile information</li>
					<li>Published blog posts</li>
					<li>Comments and interactions</li>
					<li>Account preferences</li>
				</ul>
			</div>
			<p class="message">
				We're sorry to see you go, but we respect your decision. Thank you for being part of the BlogSphere community.
			</p>
			<div class="warning">
				⚠️ <strong>Important:</strong> If this deletion was made in error or you have any concerns, please contact our support team immediately. Note that account recovery may not be possible after deletion.
			</div>
			<p class="message">
				We'd love to hear your feedback about your experience on BlogSphere. Your insights help us improve our platform for the developer community.
			</p>
			<p class="message" style="margin-top: 30px; color: #18181b;">
				We hope to see you again in the future!<br><br>
				<strong>Best regards,</strong><br>
				The BlogSphere Team
			</p>
		</div>
		<div class="footer">
			<p>This is an automated confirmation email for account deletion.</p>
			<p>© 2026 BlogSphere. All rights reserved.</p>
		</div>
	</div>
</body>
</html>`
