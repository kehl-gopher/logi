package mailer

type Topic string

const (
	SendVerificationEmail Topic = "send-verification-email"
	SendWelcomeemail      Topic = "send-welcome-email"
)

func (t Topic) GetValue() string {
	switch t {
	case SendVerificationEmail:
		return "send-verification-email"
	case SendWelcomeemail:
		return "send-welcome-email"
	default:
		return ""
	}
}

type SendWelcomeEmail struct {
	Email string `json:"email"`
}


