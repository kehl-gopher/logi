package mailer

import (
	"bytes"
	"embed"
	"errors"
	"fmt"
	"html/template"
	"log"
	"net"
	"path/filepath"
	"strings"
	"time"

	"github.com/go-mail/mail/v2"
	"github.com/kehl-gopher/logi/internal/config"
	"github.com/kehl-gopher/logi/internal/utils"
)

//go:embed "templates"
var templateFS embed.FS

type Mailer struct {
	dialer *mail.Dialer
	sender string
	Conf   *config.Config
	Logs   *utils.Log
}

type Email struct {
	Recipient string `json:"recipient"`
	Subject   string `json:"subject"`
	Att       string
	AttName   string
	data      map[string]interface{}
}

func NewMailer(host string, port int, username, password, sender string, conf *config.Config, log *utils.Log) Mailer {
	dialer := mail.NewDialer(host, port, username, password)
	dialer.RetryFailure = true
	dialer.Timeout = 5 * time.Second
	return Mailer{dialer: dialer, sender: sender, Conf: conf, Logs: log}
}

func NewEmail(to string, body []byte, subject string, att string, attName string, data map[string]interface{}) Email {
	return Email{Recipient: to, Subject: subject, Att: att, AttName: attName, data: data}
}

func (m Mailer) SendEmail(recipient, temp string, e Email) error {
	var body bytes.Buffer

	e.data = addDataTemplate(e.data, m.Conf)
	tmpl, err := template.New("email").ParseFS(templateFS, filepath.Join("templates", temp))
	if err != nil {
		log.Println(err.Error())
		return err
	}
	err = tmpl.Execute(&body, e.data)
	if err != nil {
		return err
	}

	if err := m.sendEmail(body.String(), e); err != nil {

	}
	return nil
}

func (m Mailer) sendEmail(body string, e Email) error {
	msg := mail.NewMessage()
	msg.SetHeader("Subject", e.Subject)
	msg.SetHeader("From", m.sender)
	msg.SetHeader("Content-Type", "text/html; charset=UTF-8")
	msg.SetHeader("To", e.Recipient)
	msg.SetBody("text/html", body)

	const maxRetries = 5
	var lastErr error

	for attempt := 1; attempt <= maxRetries; attempt++ {
		err := m.dialer.DialAndSend(msg)
		if err == nil {
			return nil
		}

		if strings.Contains(err.Error(), "535") || strings.Contains(err.Error(), "550") {
			return fmt.Errorf("unexpected email error: %v", err)
		}

		if isNetworkError(err) {
			lastErr = err
			time.Sleep(time.Second * time.Duration(attempt))
			continue
		}

		return err
	}
	defer msg.Reset()

	return lastErr
}

func isNetworkError(err error) bool {
	var netErr net.Error
	var dnsErr *net.DNSError

	if errors.As(err, &netErr) {
		return true
	}

	if errors.As(err, &dnsErr) {
		return true
	}

	if strings.Contains(err.Error(), "connection refused") ||
		strings.Contains(err.Error(), "connection reset") ||
		strings.Contains(err.Error(), "timeout") ||
		strings.Contains(err.Error(), "temporary failure") {
		return true
	}

	return false
}

func addDataTemplate(data map[string]interface{}, conf *config.Config) map[string]interface{} {
	data["frontend_url"] = conf.APP_CONFIG.FRONTEND_URL

	return data
}
