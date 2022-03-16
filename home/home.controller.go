package home

import (
	"net/http"

	"gitlab.finema.co/finema/etda/key-repository-api/requests"
	"gitlab.finema.co/finema/etda/key-repository-api/services"
	core "ssi-gitlab.teda.th/ssi/core"
	"ssi-gitlab.teda.th/ssi/core/utils"
)

type HomeController struct{}

func (n *HomeController) Get(c core.IHTTPContext) error {
	homeSvc := services.NewHomeService(c)
	response, ierr := homeSvc.Home()
	if ierr != nil {
		return c.JSON(ierr.GetStatus(), ierr.JSON())
	}

	return c.JSON(http.StatusOK, response)
}

func (n *HomeController) Store(c core.IHTTPContext) error {
	input := &requests.KeyStore{}
	if err := c.BindWithValidate(input); err != nil {
		return c.JSON(err.GetStatus(), err.JSON())
	}

	keySvc := services.NewKeyService(c, services.NewHSMService(c))
	key, ierr := keySvc.Store(&services.KeyStorePayload{
		PublicKey:  utils.GetString(input.PublicKey),
		PrivateKey: utils.GetString(input.PrivateKey),
		KeyType:    utils.GetString(input.KeyType),
	})
	if ierr != nil {
		return c.JSON(ierr.GetStatus(), ierr.JSON())
	}

	return c.JSON(http.StatusCreated, key)
}

func (n *HomeController) Generate(c core.IHTTPContext) error {
	keySvc := services.NewKeyService(c, services.NewHSMService(c))
	key, ierr := keySvc.Generate()
	if ierr != nil {
		return c.JSON(ierr.GetStatus(), ierr.JSON())
	}

	return c.JSON(http.StatusCreated, key)
}

func (n *HomeController) GenerateRSA(c core.IHTTPContext) error {
	keySvc := services.NewKeyService(c, services.NewHSMService(c))
	key, ierr := keySvc.GenerateRSA()
	if ierr != nil {
		return c.JSON(ierr.GetStatus(), ierr.JSON())
	}

	return c.JSON(http.StatusCreated, key)
}

func (n *HomeController) Sign(c core.IHTTPContext) error {
	input := &requests.KeySign{}
	if err := c.BindWithValidate(input); err != nil {
		return c.JSON(err.GetStatus(), err.JSON())
	}

	keySvc := services.NewKeyService(c, services.NewHSMService(c))
	key, ierr := keySvc.Sign(utils.GetString(input.ID), utils.GetString(input.Message))
	if ierr != nil {
		return c.JSON(ierr.GetStatus(), ierr.JSON())
	}

	return c.JSON(http.StatusOK, core.Map{
		"signature": key,
		"message":   utils.GetString(input.Message),
	})
}
