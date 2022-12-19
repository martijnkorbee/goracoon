package mailer

import (
	"bytes"
	"errors"
	"fmt"
	"html/template"
	"os"
	"time"

	mail "github.com/xhit/go-simple-mail/v2"
)

type Mail struct {
	Domain      string
	Templates   string
	Host        string
	Port        int
	Username    string
	Password    string
	Encryption  string
	FromAddress string
	FromName    string
	Jobs        chan Message
	Results     chan Result
	API         string
	API_KEY     string
	API_URL     string
}

type Message struct {
	From        string
	FromName    string
	To          []string
	CC          []string
	BCC         []string
	Subject     string
	Template    string
	Attachments []string
	Data        interface{}
}

type Result struct {
	Success bool
	Error   error
}

func (m *Mail) ListenForMail() {
	for {
		msg := <-m.Jobs
		err := m.Send(msg)
		if err != nil {
			m.Results <- Result{false, err}
		} else {
			m.Results <- Result{true, nil}
		}
	}
}

func (m *Mail) Send(msg Message) error {
	switch os.Getenv("MAILER") {
	case "SMTP":
		return m.SendSMTPMessage(msg)
	default:
		// no mailer specified
		return errors.New("none or invalid mailer specified in .env file")
	}
}

func (m *Mail) SendSMTPMessage(msg Message) error {
	formattedMsg, err := m.buildHTMLMessage(msg)
	if err != nil {
		return err
	}

	plainTextMsg, err := m.buildPlainTextMessage(msg)
	if err != nil {
		return err
	}

	server := mail.NewSMTPClient()

	server.Host = m.Host
	server.Port = m.Port
	server.Username = m.Username
	server.Password = m.Password
	// set encryption
	switch os.Getenv("SMTP_ENCRYPTION") {
	case "TLS":
		server.Encryption = mail.EncryptionSTARTTLS
	case "SSL":
		server.Encryption = mail.EncryptionSSLTLS
	case "none":
		server.Encryption = mail.EncryptionNone
	default:
		server.Encryption = mail.EncryptionSTARTTLS
	}
	// set authentication
	switch os.Getenv("SMTP_AUTH_METHOD") {
	case "login":
		server.Authentication = mail.AuthLogin
	case "CRAMMD5":
		server.Authentication = mail.AuthCRAMMD5
	default:
		server.Authentication = mail.AuthPlain
	}

	server.KeepAlive = false

	server.ConnectTimeout = 15 * time.Second
	server.SendTimeout = 15 * time.Second

	smtpClient, err := server.Connect()
	if err != nil {
		return err
	}

	email := mail.NewMSG()
	email.SetFrom(m.FromAddress).
		AddTo(msg.To...).
		SetSubject("test email").
		AddCc(msg.CC...).
		AddBcc(msg.BCC...).
		SetBody(mail.TextHTML, formattedMsg).
		AddAlternative(mail.TextPlain, plainTextMsg)

	// TODO: add attachments

	err = email.Send(smtpClient)
	if err != nil {
		return err
	}

	return nil
}

func (m *Mail) buildHTMLMessage(msg Message) (string, error) {

	templateToRender := fmt.Sprintf("%s/%s.html.tmpl", m.Templates, msg.Template)

	// load a html template
	t, err := template.New("email-html").ParseFiles(templateToRender)
	if err != nil {
		return "", err
	}

	// execute and store in a io buffer
	var tpl bytes.Buffer
	err = t.ExecuteTemplate(&tpl, "body", msg.Data)
	if err != nil {
		return "", err
	}

	return tpl.String(), nil
}

func (m *Mail) buildPlainTextMessage(msg Message) (string, error) {

	templateToRender := fmt.Sprintf("%s/%s.plain.tmpl", m.Templates, msg.Template)

	// load a template
	t, err := template.New("email-plain").ParseFiles(templateToRender)
	if err != nil {
		return "", err
	}

	// execute and store in a io buffer
	var tpl bytes.Buffer
	err = t.ExecuteTemplate(&tpl, "body", msg.Data)
	if err != nil {
		return "", err
	}

	return tpl.String(), nil
}
