package mailer

import (
	"bytes"
	"fmt"
	"html/template"

	"github.com/kehl-gopher/logi/internal/config"
)

type EmailType string

var (
	VerificationEmail   EmailType
	WelcomeEmail        EmailType
	ForgotPasswordEmail EmailType
)

func (t EmailType) String() string {
	switch t {
	case VerificationEmail:
		return "send-verification-email"
	case WelcomeEmail:
		return "send-welcome-email"
	case ForgotPasswordEmail:
		return "forgot-password"
	default:
		return ""
	}
}

type EmailJOB struct {
	Type EmailType
	To   string
	Data map[string]interface{} `json:"data"`
}

func (job EmailJOB) HandleEmailJob() (string, string, error) {
	body, err := job.parseEmailTemplate()
	if err != nil {
		return "", "", err
	}

	subject := map[EmailType]string{
		WelcomeEmail:        "Welcome to Logi 🎉",
		VerificationEmail:   "Verify your email ✅",
		ForgotPasswordEmail: "Reset your password 🔐",
	}[job.Type]

	return subject, body, err
}

func (job EmailJOB) parseEmailTemplate() (string, error) {

	var templ string

	switch job.Type {
	case WelcomeEmail:
		templ = "templates/welcome_email.html"
	case ForgotPasswordEmail:
		templ = "templates/forgot_password_email.html"
	case VerificationEmail:
		templ = "templates/verification_email.html"
	default:
		return "", fmt.Errorf("email template not found...")
	}

	tmpl, err := template.New("email").ParseFS(templateFS, templ)

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
