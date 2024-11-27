package utils

import (
	"gopkg.in/gomail.v2"
)

func SendEmail(to, subject, body string) error {
	m := gomail.NewMessage()
	m.SetHeader("From", "your_email@example.com")
	m.SetHeader("To", to)
	m.SetHeader("Subject", subject)
	m.SetBody("text/plain", body)

	d := gomail.NewDialer("smtp.example.com", 587, "your_email@example.com", "your_password")

	return d.DialAndSend(m)
}