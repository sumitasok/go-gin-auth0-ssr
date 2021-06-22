package main

import (
	"github.com/asteriaaerospace/back-office/handler"
	logConfig "github.com/asteriaaerospace/back-office/library/log"
	ginTemplate "github.com/foolin/gin-template"
	gin "github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"time"
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

	home := handler.Home{}

	r.GET("/root", home.Root)

	r.Run(serverPort)
	
}
