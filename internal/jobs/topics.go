package jobs

const (
	SendWelcomeEmail        = "send-welcome-email"
	SendForgotPasswordEmail = "send-forgot-password-email"
)

type SendEmail struct {
	RecipientEmail string `json:"recipient"`
	RecipientName  string `json:"name"`
	Body           []byte `json:"body"`
}
