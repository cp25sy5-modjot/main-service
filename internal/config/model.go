package config

type Config struct {
	PostgreSQL *PostgreSQL
	App        *Fiber
	Auth       *Auth
	Google     *Google
	AIService  *AIService
}

type Fiber struct {
	Host string
	Port string
}

// Database
type PostgreSQL struct {
	Host     string
	Port     string
	Protocol string
	Username string
	Password string
	Database string
	SSLMode  string
}

type Auth struct {
	AccessTokenSecret  string
	RefreshTokenSecret string
	AccessTokenTTL     string
	RefreshTokenTTL    string
	Issuer             string
}

type Google struct {
	ClientID string
}

type AIService struct {
	Url string
}