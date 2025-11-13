package auth

import (
	"github.com/Hajdudev/invoice-flow/internal/env"
	"github.com/gorilla/sessions"
	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
	"github.com/markbates/goth/providers/google"
)

const (
	maxAge = 86400 * 30
)

func NewOauth() {
	googleClientID := env.MustGetString("GOOGLE_CLIENT_ID")
	googleClientSecret := env.MustGetString("GOOGLE_CLIENT_SECRET")
	key := env.GetString("KEY", "sksdfjl")
	isProd, err := env.IsProduction()
	if err != nil {
		panic(err)
	}

	store := sessions.NewCookieStore([]byte(key))
	store.MaxAge(maxAge)

	store.Options.Path = "/"
	store.Options.HttpOnly = true
	store.Options.Secure = isProd

	gothic.Store = store
	goth.UseProviders(
		google.New(googleClientID, googleClientSecret, "http://localhost:3000/auth/google/callback"),
	)
}
