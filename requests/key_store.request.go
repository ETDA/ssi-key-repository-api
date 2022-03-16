package requests

import (
	"fmt"
	"gitlab.finema.co/finema/etda/key-repository-api/consts"
	core "ssi-gitlab.teda.th/ssi/core"
)

type KeyStore struct {
	core.BaseValidator
	PublicKey  *string `json:"public_key"`
	PrivateKey *string `json:"private_key"`
	KeyType    *string `json:"key_type"`
}

func (r KeyStore) Valid(ctx core.IContext) core.IError {
	r.Must(r.IsStrRequired(r.PublicKey, "public_key"))
	r.Must(r.IsStrRequired(r.PrivateKey, "private_key"))
	r.Must(r.IsStrIn(r.KeyType, fmt.Sprintf("%s|%s", consts.KeyTypeECDSA, consts.KeyTypeRSA), "key_type"))
	r.Must(r.IsStrRequired(r.KeyType, "key_type"))

	return r.Error()
}
