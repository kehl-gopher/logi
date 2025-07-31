package mailer

import (
	"embed"
	"errors"
	"fmt"
	"net"
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
	Conf   *config.AppConfig
	Logs   *utils.Log
}

type Email struct {
	Recipient string `json:"recipient"`
	Subject   string `json:"subject"`
	Att       string
	AttName   string
	data      map[string]interface{}
}

func NewMailer(host string, port int, username, password, sender string, conf *config.AppConfig, log *utils.Log) Mailer {
	dialer := mail.NewDialer(host, port, username, password)
	dialer.RetryFailure = true
	dialer.Timeout = 5 * time.Second
	return Mailer{dialer: dialer, sender: sender, Conf: conf, Logs: log}
}

func NewEmail(to string, subject string, att string, attName string, data map[string]interface{}) Email {
	return Email{Recipient: to, Subject: subject, Att: att, AttName: attName, data: data}
}

func email(conf *config.AppConfig, e EmailJOB, lg *utils.Log) error {
	port, err := utils.PortResolver(conf.SMTP_PORT)
	if err != nil {
		return err
	}

	m := NewMailer(conf.SMTP_HOST, port, conf.SMTP_USERNAME, conf.SMTP_PASSWORD, conf.SMTP_USERNAME, conf, lg)

	body, subject, err := e.HandleEmailJob()

	if err != nil {
		return err
	}

	if err := m.sendEmail(body, subject, e); err != nil {
		utils.PrintLog(lg, fmt.Sprintf("failed to send email: %v", err), utils.ErrorLevel)
		return err
	}

	return nil
}

func (m Mailer) sendEmail(body string, subject string, e EmailJOB) error {
	msg := mail.NewMessage()
	msg.SetHeader("Subject", subject)
	msg.SetAddressHeader("From", m.sender, "Logi Team") // Add name
	msg.SetBody("text/html", body)
	msg.SetHeader("To", e.To)

	const maxRetries = 3
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

	if lastErr != nil {
		return lastErr
	}
	defer msg.Reset()
	return nil
}

func (e *EmailJOB) SendWelcomeEmails(conf *config.AppConfig, lg *utils.Log) error {
	if err := email(conf, *e, lg); err != nil {
		return err
	}
	return nil
}

func (e *EmailJOB) SendVerificationMail(conf *config.AppConfig, lg *utils.Log) error {
	if err := email(conf, *e, lg); err != nil {
		return err
	}
	return nil
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

func addDataTemplate(data map[string]interface{}, conf *config.AppConfig) map[string]interface{} {
	data["frontend_url"] = conf.FRONTEND_URL
	return data
}
