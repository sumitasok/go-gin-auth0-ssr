package handler

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"net/http"
)

type Authentication struct {
}

func (a Authentication) Login(c *gin.Context) {
	session := sessions.Default(c)
	session.Set("id", 12090292)
	session.Set("email", "test@gmail.com")
	_ = session.Save() // handle error
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
	ctx.JSON(http.StatusOK, gin.H{
		"message": "callback received",
	})
}
