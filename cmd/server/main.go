package main

import (
	"encoding/gob"
	"github.com/sumitasok/back-office/library/middlewares"

	//"github.com/sumitasok/back-office/library/middlewares"
	"github.com/gin-contrib/sessions"
	"time"

	ginTemplate "github.com/foolin/gin-template"
	gin "github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"github.com/sumitasok/back-office/handler"
	logConfig "github.com/sumitasok/back-office/library/log"

	"github.com/gin-contrib/sessions/redis"

	"github.com/joho/godotenv"
)

var r = gin.New()
var serverPort = ":7271"
var authSessionName = "authsession" // os.Getenv("AUTH_SESSION_NAME")

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Print(err.Error())
	}

	gob.Register(map[string]interface{}{})

	r.HTMLRender = ginTemplate.New(ginTemplate.TemplateConfig{
		Root:      "templates",
		Extension: ".gohtml",
		Master:    "layouts/main",
		Partials:  []string{},
	})
}

func main() {

	startAt := time.Now()
	log.Infoln("server starting...")
	r.Use(logConfig.WithLogrus())
	r.Use(gin.Recovery())
	log.Infoln("middlewares set in...", time.Since(startAt))

	store, _ := redis.NewStore(10, "tcp", "localhost:6379", "", []byte("secret"))
	r.Use(sessions.Sessions(authSessionName, store))

	home := handler.Home{}
	auth := handler.Authentication{
		LandingPage: "/manage/root",
	}

	r.GET("/root", home.Root)
	r.GET("/login", auth.Login)
	r.GET("/logout", auth.Logout)
	r.GET("/callback", auth.Callback)

	authR := r.Group("/manage")
	authR.Use(middlewares.Auth0Authentication())

	authR.GET("/root", home.Root)

	_ = r.Run(serverPort) // handle error.

}
