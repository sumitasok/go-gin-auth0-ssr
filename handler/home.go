package handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type Home struct {
}

func (h Home) Root(ctx *gin.Context) {
	ctx.HTML(http.StatusOK, "views/home/root",
		gin.H{
			"Name": "Golang",
		})
}
