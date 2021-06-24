package handler

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"github.com/asteriaaerospace/back-office/library/auth"
	"github.com/coreos/go-oidc"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"net/http"
	"os"
)

type Authentication struct {
}

func (a Authentication) Login(c *gin.Context) {
	// session sample
	session := sessions.Default(c)
	//session.Set("_id", 12090292)
	//session.Set("_email", "test@gmail.com")
	//_ = session.Save() // handle error

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

	c.JSON(http.StatusOK, gin.H{
		"message": "User Sign In successfully",
	})
}

func (a Authentication) Logout(c *gin.Context) {
	session := sessions.Default(c)
	session.Clear()
	_ = session.Save() // handle error
	c.JSON(http.StatusOK, gin.H{
		"message": "User Sign out successfully",
	})
}

func (a Authentication) Callback(ctx *gin.Context) {
	session := sessions.Default(ctx)
	sessionState, ok := session.Get("state").(string)

	session.Set("KEY", "SESSION")

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

	log.Info("code: " + ctx.Query("code"))
	log.Info("authenticator.Config: ", authenticator.Config)
	log.Info("authenticator.Config.Provider: ", authenticator.Config.Endpoint.TokenURL)
	token, err := authenticator.Config.Exchange(context.TODO(), ctx.Query("code"))
	if err != nil {
		log.Printf("no token found: %v", err)
		log.Error("2: ", err.Error())
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

	// if all pass

	ctx.Redirect(http.StatusSeeOther, "/manage/root")
}
