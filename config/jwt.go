package config

import (
	"log"
	"os"
)

var JWTKey []byte

func InitJWT() {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		log.Fatal("JWT_SECRET not set in environment")
	}
	JWTKey = []byte(secret)
}
