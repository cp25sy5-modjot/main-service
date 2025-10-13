package config

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
	AccessTokenExpireMinutes  int64
	RefreshTokenExpireMinutes int64
	SecretKey                 string
	Issuer                    string
}

type Google struct {
	ClientID     string
	ClientSecret string
	RedirectURL  string
}
