package handler

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"github.com/coreos/go-oidc"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"github.com/sumitasok/back-office/library/auth"
	"net/http"
	"net/url"
	"os"
)

type Authentication struct {
	// LandingPage - where the user is redirected once logged in; if returnTo is not available.
	LandingPage string
}

func (a Authentication) Login(c *gin.Context) {
	session := sessions.Default(c)

	// auth0 login
	b := make([]byte, 32)
	_, err := rand.Read(b)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		return
	}
	state := base64.StdEncoding.EncodeToString(b)

	session.Set("state", state)
	err = session.Save()

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		return
	}

	authenticator, err := auth.NewAuthenticator()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		return
	}

	c.Redirect(http.StatusTemporaryRedirect, authenticator.Config.AuthCodeURL(state))
}

func (a Authentication) Logout(ctx *gin.Context) {
	// clear up the session set.
	session := sessions.Default(ctx)
	session.Clear()
	err := session.Save()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		return
	}

	domain := os.Getenv("AUTH0_DOMAIN")

	logoutUrl, err := url.Parse("https://" + domain)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		return
	}

	logoutUrl.Path += "/v2/logout"
	parameters := url.Values{}

	var scheme string
	if ctx.Request.TLS == nil {
		scheme = "http"
	} else {
		scheme = "https"
	}

	returnTo, err := url.Parse(scheme + "://" + ctx.Request.Host)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		return
	}
	parameters.Add("returnTo", returnTo.String())
	parameters.Add("client_id", os.Getenv("AUTH0_CLIENT_ID"))
	logoutUrl.RawQuery = parameters.Encode()

	ctx.Redirect(http.StatusTemporaryRedirect, logoutUrl.String())
}

func (a Authentication) Callback(ctx *gin.Context) {
	session := sessions.Default(ctx)
	sessionState, ok := session.Get("state").(string)

	// session.State
	if !ok {
		log.Error("sessionState not found")
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "sessionState not found",
		})
	}

	if ctx.Query("state") != sessionState {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid state parameter",
		})
		return
	}

	// Authenticator
	// https://manage.auth0.com/dashboard/us/cartis/applications/RefYtdrjJwKUarNcbPSAg8fajE2kqVHw/quickstart
	authenticator, err := auth.NewAuthenticator()
	if err != nil {
		log.Error("1: ", err.Error())
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		return
	}

	token, err := authenticator.Config.Exchange(context.TODO(), ctx.Query("code"))
	if err != nil {
		log.Errorf("no token found: %v", err)
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"message": err.Error(),
		})
		return
	}

	rawIDToken, ok := token.Extra("id_token").(string)
	if !ok {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "No id_token field in oauth2 token.",
		})
		return
	}

	oidcConfig := &oidc.Config{
		ClientID: os.Getenv("AUTH0_CLIENT_ID"),
	}

	idToken, err := authenticator.Provider.Verifier(oidcConfig).Verify(context.TODO(), rawIDToken)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to verify ID Token: " + err.Error(),
		})
		return
	}

	var profile map[string]interface{}
	if err := idToken.Claims(&profile); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		return
	}

	session.Set("id_token", rawIDToken)
	session.Set("access_token", token.AccessToken)
	session.Set("profile", profile)

	err = session.Save()

	ctx.Redirect(http.StatusSeeOther, a.LandingPage)
}
