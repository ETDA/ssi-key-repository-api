package services

import (
	"crypto/x509"
	"errors"

	"gitlab.finema.co/finema/etda/key-repository-api/consts"
	"gitlab.finema.co/finema/etda/key-repository-api/emsgs"
	"gitlab.finema.co/finema/etda/key-repository-api/models"
	core "ssi-gitlab.teda.th/ssi/core"
	"ssi-gitlab.teda.th/ssi/core/errmsgs"
	"ssi-gitlab.teda.th/ssi/core/utils"
	"gorm.io/gorm"
)

type KeyGeneratePayload struct {
}

type KeySignPayload struct {
	ID      string
	Message string
}

type KeyStorePayload struct {
	PublicKey  string
	PrivateKey string
	KeyType    string
}

type IKeyService interface {
	Find(id string) (*models.Key, core.IError)
	Store(payload *KeyStorePayload) (*models.Key, core.IError)
	Generate() (*models.Key, core.IError)
	GenerateRSA() (*models.Key, core.IError)
	Sign(id string, message string) (string, core.IError)
}
type keyService struct {
	ctx        core.IContext
	hsmService IHSMService
}

func NewKeyService(ctx core.IContext, hsmService IHSMService) IKeyService {
	return &keyService{
		ctx:        ctx,
		hsmService: hsmService,
	}
}

func (s keyService) Find(id string) (*models.Key, core.IError) {
	key := &models.Key{}
	err := s.ctx.DB().First(&key, "id = ?", id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, s.ctx.NewError(err, emsgs.KeyNotFoundError)
	}
	if err != nil {
		return nil, s.ctx.NewError(err, errmsgs.DBError)
	}

	return key, nil
}

func (s keyService) Generate() (*models.Key, core.IError) {
	kp, err := utils.GenerateKeyPair()
	if err != nil {
		return nil, s.ctx.NewError(err, emsgs.GenerateKeyError)
	}

	return s.Store(&KeyStorePayload{
		PublicKey:  kp.PublicKeyPem,
		PrivateKey: kp.PrivateKeyPem,
		KeyType:    string(consts.KeyTypeECDSA),
	})
}

func (s keyService) GenerateRSA() (*models.Key, core.IError) {
	kp, err := utils.GenerateKeyPairWithOption(&utils.GenerateKeyPairOption{
		Algorithm: x509.SHA256WithRSA,
		KeySize:   utils.EncryptRSA2048Bits,
	})

	if err != nil {
		return nil, s.ctx.NewError(err, emsgs.GenerateKeyError)
	}

	rsaKeyPair, ok := kp.(*utils.RSAKeyPair)
	if !ok {
		return nil, s.ctx.NewError(emsgs.GenerateKeyError, emsgs.GenerateKeyError)
	}

	return s.Store(&KeyStorePayload{
		PublicKey:  rsaKeyPair.PublicKeyPem,
		PrivateKey: rsaKeyPair.PrivateKeyPem,
		KeyType:    string(consts.KeyTypeRSA),
	})
}

func (s keyService) Sign(id string, message string) (string, core.IError) {
	key, ierr := s.Find(id)
	if ierr != nil {
		return "", s.ctx.NewError(ierr, ierr)
	}

	decryptedPrivateKey, ierr := s.hsmService.Decrypt(key.PrivateKeyEncrypted)
	if ierr != nil {
		return "", s.ctx.NewError(ierr, ierr)
	}

	var signature string
	if key.Type == string(consts.KeyTypeECDSA) {
		privateKey, err := utils.LoadPrivateKey(decryptedPrivateKey)
		if err != nil {
			return "", s.ctx.NewError(ierr, errmsgs.InternalServerError)
		}
		signature, err = utils.SignMessage(privateKey, message)
		if err != nil {
			return "", s.ctx.NewError(ierr, errmsgs.InternalServerError)
		}
	} else if key.Type == string(consts.KeyTypeRSA) {
		privateKey, err := utils.LoadRSAPrivateKey(decryptedPrivateKey)
		if err != nil {
			return "", s.ctx.NewError(ierr, errmsgs.InternalServerError)
		}
		signature, err = utils.SignMessageWithOption(privateKey, message, &utils.SignMessageOption{
			Algorithm: x509.SHA256WithRSA,
		})
		if err != nil {
			return "", s.ctx.NewError(ierr, errmsgs.InternalServerError)
		}
	} else {
		return "", s.ctx.NewError(emsgs.UnsupportedSigningAlgorithm, emsgs.UnsupportedSigningAlgorithm)
	}

	return signature, nil
}

func (s keyService) Store(payload *KeyStorePayload) (*models.Key, core.IError) {
	encryptedPrivateKey, ierr := s.hsmService.Encrypt(payload.PrivateKey)
	if ierr != nil {
		return nil, s.ctx.NewError(ierr, ierr)
	}

	key := models.NewKey(payload.PublicKey, encryptedPrivateKey, payload.KeyType)
	err := s.ctx.DB().Create(key).Error
	if err != nil {
		return nil, s.ctx.NewError(err, errmsgs.DBError)
	}

	return s.Find(key.ID)
}
