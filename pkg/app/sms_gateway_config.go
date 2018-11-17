package app

// SmsGatewayConfig stores data for authenticating
// on SMS service providers
type SmsGatewayConfig struct {
	SenderID string
	Username string
	Password string
}

// NewSmsGatewayConfig creates and returns a pointer to a SmsGatewayConfig
func NewSmsGatewayConfig(senderID, username, password string) *SmsGatewayConfig {
	return &SmsGatewayConfig{
		SenderID: senderID,
		Username: username,
		Password: password,
	}
}
