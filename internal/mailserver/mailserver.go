package mailserver

import (
	"gopkg.in/gomail.v2"
)

type Config struct {
	Host     string `json:"host"`
	Port     int    `json:"port"`
	Username string `json:"username"`
	Password string `json:"password"`
}

type MailServer struct {
	dialer *gomail.Dialer

	defFrom     string
	defFromName string
}

func NewMailServer(config *Config, defFrom, defFromName string) (*MailServer, error) {
	ms := new(MailServer)
	ms.dialer = gomail.NewPlainDialer(config.Host, config.Port, config.Username, config.Password)

	closer, err := ms.dialer.Dial()
	if err != nil {
		return nil, err
	}

	defer closer.Close()

	ms.defFrom = defFrom
	ms.defFromName = defFromName

	return ms, nil
}

func (ms *MailServer) SendMailRaw(msg *gomail.Message) error {
	return ms.dialer.DialAndSend(msg)
}

func (ms *MailServer) SendMail(from, fromName, to, subject, body, bodyType string) error {
	if from == "" {
		from = ms.defFrom
	}

	if fromName == "" {
		fromName = ms.defFromName
	}

	msg := gomail.NewMessage()
	msg.SetAddressHeader("From", from, fromName)
	msg.SetAddressHeader("To", to, "")
	msg.SetHeader("Subject", subject)
	msg.SetBody(bodyType, body)
	return ms.SendMailRaw(msg)
}

func (ms *MailServer) SendMailFromDef(to, subject, body, bodyType string) error {
	return ms.SendMail("", "", to, subject, body, bodyType)
}
