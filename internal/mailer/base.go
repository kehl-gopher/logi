package mailer

import (
	"bytes"
	"fmt"

	"html/template"

	"github.com/kehl-gopher/logi/internal/config"
)

type EmailType string

const (
	VerificationEmail   EmailType = "send-verification-email"
	WelcomeEmail        EmailType = "send-welcome-email"
	ForgotPasswordEmail EmailType = "send-forgot-password-mail"
)

type EmailJOB struct {
	Type EmailType              `json:"type"`
	To   string                 `json:"to"`
	Data map[string]interface{} `json:"data"`
}

func NewEmailJob(typed EmailType, to string, data map[string]interface{}) EmailJOB {
	if data == nil {
		data = make(map[string]interface{})
	}
	return EmailJOB{To: to, Type: typed, Data: data}
}
func (job *EmailJOB) HandleEmailJob() (string, string, error) {
	body, err := job.parseEmailTemplate()
	if err != nil {
		return "", "", err
	}

	subject := map[EmailType]string{
		WelcomeEmail:        "Welcome to Logi 🎉",
		VerificationEmail:   "Verify your email ✅",
		ForgotPasswordEmail: "Reset your password 🔐",
	}[job.Type]

	return body, subject, err
}

func (job *EmailJOB) parseEmailTemplate() (string, error) {
	var templ string

	switch job.Type {
	case WelcomeEmail:
		templ = "templates/welcome_email.html"
	case ForgotPasswordEmail:
		templ = "templates/forgot_password_email.html"
	case VerificationEmail:
		templ = "templates/verification_mail.html"
	default:
		return "", fmt.Errorf("email template not found")
	}

	tmpl, err := template.ParseFS(templateFS, templ)

	if err != nil {
		return "", err
	}

	var body bytes.Buffer

	data := addDataTemplate(job.Data, &config.AppConfig{})
	err = tmpl.Execute(&body, data)
	if err != nil {
		return "", err
	}
	return body.String(), nil
}
