package email

import (
	"net/smtp"

	"github.com/jordan-wright/email"
)

func SendEmail(mailAddr, pwd, addr string) {
	e := email.NewEmail()
	e.From = mailAddr
	e.To = []string{"kirill@zemskoff.su"}
	e.Subject = "Verify"
	e.HTML = []byte("<a>http://somedomain.ru</a>")
	e.Send("smtp.gmail.com:587", smtp.PlainAuth("", mailAddr, pwd, addr))
	// fmt.Println("email send")
}
