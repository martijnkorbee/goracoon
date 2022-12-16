package mailer

import (
	"bytes"
	"errors"
	"fmt"
	"html/template"
	"os"

	"github.com/ainsleyclark/go-mail/drivers"
	"github.com/ainsleyclark/go-mail/mail"
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

	cfg := mail.Config{
		URL:         m.Host,
		FromAddress: m.FromAddress,
		FromName:    m.FromName,
		Password:    m.Password,
		Port:        m.Port,
	}

	mailer, err := drivers.NewSMTP(cfg)
	if err != nil {
		return err
	}

	tx := &mail.Transmission{
		Recipients: msg.To,
		Subject:    msg.Subject,
		HTML:       formattedMsg,
		PlainText:  plainTextMsg,
	}

	// TODO: add attachments

	result, err := mailer.Send(tx)
	if err != nil {
		return err
	}

	fmt.Printf("%+v\n", result)

	return nil
}

func (m *Mail) buildHTMLMessage(msg Message) (string, error) {

	templateToRender := fmt.Sprintf("%s/%s/.html.tmpl", m.Templates, msg.Template)

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

	templateToRender := fmt.Sprintf("%s/%s/.plain.tmpl", m.Templates, msg.Template)

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
