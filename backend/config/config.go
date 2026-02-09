package config

import (
	"log"
	"os"

	"github.com/OderoCeasar/joho/godotenv"
)


type Config struct {
	Server		ServerConfig
	database	DatabaseConfig
	JWT			JWTConfig
	Mpesa		MpesaConfig
	RADIUS		RADIUSConfig
	MikroTik	MikroTikConfig
}


type ServerConfig struct {
	Port		string
	GinMode 	string
	FrontendURL	string
}


type DatabaseConfig struct {
	Host		string
	Port 		string
	User 		string
	Password	string
	DBName 		string
	SSLMode		string
}


type JWTConfig struct {
	Secret 		string
}


type MpesaConfig struct {
	ConsumerKey 	string
	ConsumerSecret	string
	PassKey			string
	ShortCode		string
	InitiatorName	string
	InitiatorPassword	string
	Environment 	string
	CallbackURL 	string
}


type RADIUSConfig struct {
	Server		string
	Secret		string
	AuthPort	string
	AcctPort	string
}


type MikroTikConfig struct {
	Host		string
	Username	string
	Password	string
	Port 		string
}


func Load() *Config {
	if err := godotenv.Load(); err != nil {
		log.Println("No env files found")
	}

	return &Config{
		Server: ServerConfig{
			Port:		getEnv("PORT", "8080"),
			GinMode:    getEnv("GIN_MODE", "debug"),
			FrontendURL: getEnv("FRONTEND_URL", "http://localhost:5173"),
		},
		Database: DatabaseConfig{
			Host:		getEnv("DB_HOST", "localhost"),
			Port:		getEnv("DB_PORT", "5432"),
			User:		getEnv("DB_USER", "postgres"),
			Password:   getEnv("DB_PASSWORD", "(123@Ceasar!)"),
			DBName:     getEnv("DB_NAME", "wifi_billing"),
			SSLMode:    getEnv("DB_SSLMODE", "disable"),	 
		},
		JWT: JWTConfig{
			Secret: 	getEnv("JWT_SECRET", "secret_key"),
		},
		Mpesa:	MpesaConfig{
			ConsumerKey:		getEnv("MPESA_CONSUMER_KEY", ""),
			ConsumerSecret:		getEnv("MPESA_CONSUMER_SECRET", ""),
			PassKey:			getEnv("MPESA_PASSKEY", ""),
			ShortCode:			getEnv("MPESA_SHORTCODE", ""),
			InitiatorName:		getEnv("MPESA_INITIATOR_NAME", ""),
			InitiatorPassword:	getEnv("MPESA_INITIATOR_PASSWORD", ""),
			Environment:		getEnv("MPESA_ENVIRONMENT", "sandbox"),
			CallbackURL:		getEnv("MPESA_CALLBACK_URL", ""),
		},
		RADIUS: RADIUSConfig{
			Server:		getEnv("RADIUS_SERVER", ""),
			Secret:		getEnv("RADIUS_SECRET", ""),
			AuthPort:	getEnv("RADIUS_AUTH_PORT", ""),
			AcctPort:	getEnv("RADIUS_ACCT_PORT", ""),
		},
		MikroTik: MikroTikConfig{
			Host:		getEnv("MIKROTIK_HOST", ""),
			Username:	getEnv("MIKROTIK_USERNAME", ""),
			Password:	getEnv("MIKROTIK_PASSWORD", ""),
			Port:		getEnv("MIKROTIK_PORT", ""),
		},

	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}