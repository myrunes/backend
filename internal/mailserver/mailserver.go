package mailserver

import (
	"gopkg.in/gomail.v2"
)

// Config wraps the configuration values
// for the mail server.
type Config struct {
	Host     string `json:"host"`
	Port     int    `json:"port"`
	Username string `json:"username"`
	Password string `json:"password"`
}

// MailServer provides a connection
// to a mail-server to send e-mails.
type MailServer struct {
	dialer *gomail.Dialer

	defFrom     string
	defFromName string
}

// NewMailServer initializes a new mail server with
// the given mail server configuration, default "from"
// mail address (defFrom) and default "from" sender
// name (defFromName).
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

// SendMailRaw dials the connection to the
// mail server and sends the passed Message
// object.
func (ms *MailServer) SendMailRaw(msg *gomail.Message) error {
	return ms.dialer.DialAndSend(msg)
}

// SendMail wraps the given data from, fromName, to,
// subject, body and bodyType to a Message object
// which is then sent via SendMailRaw.
// If from and fromName is empty, the defFrom and
// defFromName will be used instead.
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

// SendMailFromDef is shorthand for SendMail
// with defFrom and defFromName as sender
// specifications.
func (ms *MailServer) SendMailFromDef(to, subject, body, bodyType string) error {
	return ms.SendMail("", "", to, subject, body, bodyType)
}
