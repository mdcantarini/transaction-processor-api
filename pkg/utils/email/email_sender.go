package email

type EmailSender interface {
	SendEmail(to, subject, body string, attachments ...string) error
}
