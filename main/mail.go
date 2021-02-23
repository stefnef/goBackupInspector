package main

import (
	"fmt"
	"goBackupInspector/config"
	"gopkg.in/gomail.v2"
	"log"
)

func sendMail(body string, summary string, conf config.Config) {
	d := gomail.NewDialer(conf.Mail.SMTPServer, conf.Mail.SMTPPort, conf.Mail.UserName, conf.Mail.Password ) //"efzoipefzsfiaefozif/QWUais")
	s, err := d.Dial()
	if err != nil {
		panic(err)
	}

	m := gomail.NewMessage()
	m.SetHeader("From", conf.Mail.UserName)
	m.SetAddressHeader("To", conf.Mail.ReceiverAdr, conf.Mail.ReceiverName)
	m.SetHeader("Subject", "Diffs in your system")
	m.SetBody("text/html", fmt.Sprintf("Hello!\n%s", body))
	m.Attach(summary)

	if err := gomail.Send(s, m); err != nil {
		log.Printf("Could not send email to %q: %v", conf.Mail.ReceiverAdr, err)
	}
	m.Reset()
}
