package services

import (
	"ssi-gitlab.teda.th/ssi/core"
)

type IHomeService interface {
	Home() (core.Map, core.IError)
}

type HomeService struct {
	ctx core.IContext
}

func NewHomeService(ctx core.IContext) IHomeService {
	return &HomeService{
		ctx: ctx,
	}
}

func (s HomeService) Home() (core.Map, core.IError) {
	return core.Map{
		"message": "Hello, I'm Home API",
	}, nil
}
