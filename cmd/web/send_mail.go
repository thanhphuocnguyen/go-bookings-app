package main

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/thanhphuocnguyen/go-bookings-app/internal/models"
	mail "github.com/xhit/go-simple-mail/v2"
)

func listenForMail() {
	go func() {
		for {
			select {
			case m := <-appConfig.MailChan:
				sendMail(m)
			}
		}
	}()
}

func sendMail(m models.MailData) {
	server := mail.NewSMTPClient()
	server.Host = "localhost"
	server.Port = 1025
	server.KeepAlive = false
	server.ConnectTimeout = 10 * time.Second
	server.SendTimeout = 10 * time.Second

	client, err := server.Connect()
	if err != nil {
		appConfig.ErrorLog.Println(err)
	}
	email := mail.NewMSG()

	email.SetFrom(m.From).AddTo(m.To).SetSubject(m.Subject)

	if m.Template == "" {
		email.SetBody(mail.TextHTML, m.Content)
	} else {
		data, err := os.ReadFile(fmt.Sprintf("./email_templates/%s", m.Template))
		if err != nil {
			appConfig.ErrorLog.Println(err)
		}
		mailTemplate := string(data)
		body := strings.Replace(mailTemplate, "[%body%]", m.Content, 1)
		email.SetBody(mail.TextHTML, body)
	}

	err = email.Send(client)
	if err != nil {
		appConfig.ErrorLog.Println(err)
	} else {
		appConfig.InfoLog.Println("Email sent")
	}
}
