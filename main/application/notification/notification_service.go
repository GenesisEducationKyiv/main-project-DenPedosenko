package notification

type AuthConfig struct {
	from     string
	password string
	smtpHost string
	smtpPort string
}

func NewAuthConfig(from, password, smtpHost, smtpPort string) AuthConfig {
	return AuthConfig{
		from:     from,
		password: password,
		smtpHost: smtpHost,
		smtpPort: smtpPort,
	}
}
