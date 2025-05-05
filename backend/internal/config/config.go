package config

import "os"

type config struct {
	SessionCookieName string
	DbHost            string
	DbPort            string
	DbUser            string
	DbPassword        string
	DbName            string
	SteamKey          string
}

var Config = config{}

func Load() {
	Config = config{
		SessionCookieName: os.Getenv("SESSION_COOKIE_NAME"),
		DbHost:            os.Getenv("DB_HOST"),
		DbPort:            os.Getenv("DB_PORT"),
		DbUser:            os.Getenv("DB_USER"),
		DbPassword:        os.Getenv("DB_PASSWORD"),
		DbName:            os.Getenv("DB_NAME"),
		SteamKey:          os.Getenv("STEAM_KEY"),
	}

}
