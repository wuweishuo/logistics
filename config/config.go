package config

type Config struct {
	Logins map[string]LoginConfig
}

type LoginConfig struct {
	Username string
	Password string
}
