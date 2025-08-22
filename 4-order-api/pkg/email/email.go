package email

import (
	"fmt"
	"net/smtp"

	"github.com/jordan-wright/email"
)

func SendEmail(from, pwd, smtpaddr, to, hash string) error {
	e := email.NewEmail()
	e.From = from
	e.To = []string{to}
	e.Subject = "Email validation"
	e.HTML = fmt.Appendf(nil, "<a href=\"http://localhost:8081/verify/%s\">Подтвердить email</a>", hash)
	err := e.Send("smtp.gmail.com:587", smtp.PlainAuth("", from, pwd, smtpaddr))
	return err
}
