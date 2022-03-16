package requests

import (
	"gitlab.finema.co/finema/etda/key-repository-api/models"
	core "ssi-gitlab.teda.th/ssi/core"
)

type KeySign struct {
	core.BaseValidator
	ID      *string `json:"id"`
	Message *string `json:"message"`
}

func (r KeySign) Valid(ctx core.IContext) core.IError {
	r.Must(r.IsStrRequired(r.ID, "id"))
	r.Must(r.IsExists(ctx, r.ID, models.Key{}.TableName(), "id", "id"))
	r.Must(r.IsStrRequired(r.Message, "message"))

	return r.Error()
}
