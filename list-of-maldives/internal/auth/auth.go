package auth

import (
	"log"
	"os"

	"github.com/gorilla/sessions"
	"github.com/joho/godotenv"
	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
	"github.com/markbates/goth/providers/google"
)

const (
	key    = "randomString"
	MaxAge = 24 * 60 * 60 * 30
	IsProd = false
)

func NewAuth() error {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	googleClientId := os.Getenv("GOOGLE_KEY")
	googleClientSecret := os.Getenv("GOOGLE_SECRET")
	backendURL := os.Getenv("BACKEND_URL")

	store := sessions.NewCookieStore([]byte(key))
	store.MaxAge(MaxAge)

	store.Options.Path = "/"
	store.Options.HttpOnly = true
	store.Options.Secure = IsProd

	gothic.Store = store
	goth.UseProviders(
		google.New(
			googleClientId,
			googleClientSecret,
			backendURL+"/auth/google/callback",
			"email", "profile",
		),
	)
	return nil
}
