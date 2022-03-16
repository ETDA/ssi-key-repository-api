package home

import (
	"github.com/labstack/echo/v4"
	core "ssi-gitlab.teda.th/ssi/core"
)

func NewHomeHTTPHandler(r *echo.Echo) {
	home := &HomeController{}

	r.GET("/", core.WithHTTPContext(home.Get))
	r.POST("/key/store", core.WithHTTPContext(home.Store))
	r.POST("/key/generate", core.WithHTTPContext(home.Generate))
	r.POST("/key/generate/rsa", core.WithHTTPContext(home.GenerateRSA))
	r.POST("/key/sign", core.WithHTTPContext(home.Sign))
}
