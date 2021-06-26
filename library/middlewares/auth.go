package middlewares

import (
	"context"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"os"

	"golang.org/x/oauth2"

	oidc "github.com/coreos/go-oidc"
)

type Auth0Authenticator struct {
	Provider *oidc.Provider
	Config   oauth2.Config
	Ctx      context.Context
}

func NewAuthenticator() (*Auth0Authenticator, error) {
	ctx := context.Background()

	provider, err := oidc.NewProvider(ctx, os.Getenv("AUTH0_ISSUER"))
	if err != nil {
		log.Printf("failed to get provider: %v", err)
		return nil, err
	}

	conf := oauth2.Config{
		ClientID:     os.Getenv("AUTH0_CLIENT_ID"),
		ClientSecret: os.Getenv("AUTH0_CLIENT_SECRET"),
		RedirectURL:  os.Getenv("AUTH0_CALLBACK_URL"),
		Endpoint:     provider.Endpoint(),
		Scopes:       []string{oidc.ScopeOpenID, "profile"},
	}

	return &Auth0Authenticator{
		Provider: provider,
		Config:   conf,
		Ctx:      ctx,
	}, nil
}

// Auth0Authentication - does authentication of the user via auth0.
func Auth0Authentication() gin.HandlerFunc {
	return func(c *gin.Context) {
		session := sessions.DefaultMany(c, os.Getenv("AUTH_SESSION_NAME"))
		sessionID := session.Get("id")
		if sessionID == nil {
			c.JSON(http.StatusNotFound, gin.H{
				"message": "unauthorized",
			})
			c.Abort()
		}
	}
}
