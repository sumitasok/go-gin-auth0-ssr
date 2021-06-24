package main

import (
	"github.com/asteriaaerospace/back-office/library/middlewares"
	"time"

	"github.com/asteriaaerospace/back-office/handler"
	logConfig "github.com/asteriaaerospace/back-office/library/log"
	ginTemplate "github.com/foolin/gin-template"
	gin "github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"

	"github.com/gin-contrib/sessions/redis"
)

var r = gin.New()
var serverPort = ":7271"

func init() {
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

	_, _ = redis.NewStore(10, "tcp", "localhost:6379", "", []byte("secret"))

	home := handler.Home{}
	auth := handler.Authentication{}

	r.GET("/root", home.Root)
	r.GET("/login", auth.Login)
	r.GET("/logout", auth.Logout)
	r.GET("/callback", auth.Callback)

	authR := r.Group("/manage")
	//authR.Use(sessions.Sessions("session-store", store))
	authR.Use(middlewares.Authentication())

	authR.GET("/root", home.Root)

	_ = r.Run(serverPort) // handle error.

}
