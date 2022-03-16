package services

import (
	"bytes"
	"errors"
	"fmt"
	"time"

	"github.com/miekg/pkcs11"
	"github.com/miekg/pkcs11/p11"
	"gitlab.finema.co/finema/etda/key-repository-api/consts"
	"gitlab.finema.co/finema/etda/key-repository-api/emsgs"
	"gitlab.finema.co/finema/etda/key-repository-api/helpers"
	core "ssi-gitlab.teda.th/ssi/core"
	"ssi-gitlab.teda.th/ssi/core/errmsgs"
)

type IHSMService interface {
	Decrypt(encryptedPrivateKey string) (string, core.IError)
	Encrypt(privateKey string) (string, core.IError)
}
type hsmService struct {
	ctx core.IContext
}

func NewHSMService(ctx core.IContext) IHSMService {
	return &hsmService{ctx: ctx}
}

func (s *hsmService) Encrypt(privateKey string) (string, core.IError) {
	message := []byte(privateKey)
	maxLength := 190 // statically set for RSA 2048 with SHA-256
	messages := make([][]byte, 0)

	i := 0
	for {
		newMassage := make([]byte, maxLength)
		if ((i + 1) * maxLength) > len(message) {
			copy(newMassage, message[i*maxLength:])
			messages = append(messages, newMassage)
			break
		}
		copy(newMassage, message[i*maxLength:((i+1)*maxLength)])
		messages = append(messages, newMassage)
		i++
	}

	cipherTexts := make([][]byte, 0)
	for _, message := range messages {
		cipherText, err := s.encrypt(message)
		if err != nil {
			return "", s.ctx.NewError(err, errmsgs.InternalServerError)
		}
		cipherTexts = append(cipherTexts, cipherText)
	}

	encryptedMessage := helpers.ByteArraySeriesToBase64StringJoined(cipherTexts, ".")

	return encryptedMessage, nil
}

func (s *hsmService) Decrypt(encryptedPrivateKey string) (string, core.IError) {
	cipherTexts, err := helpers.Base64StringJoinedToByteArraySeries(encryptedPrivateKey, ".")
	if err != nil {
		return "", s.ctx.NewError(err, errmsgs.InternalServerError)
	}
	messages := make([][]byte, 0)
	for _, cipherText := range cipherTexts {
		message, err := s.decrypt(cipherText)
		if err != nil {
			return "", s.ctx.NewError(err, errmsgs.InternalServerError)
		}
		messages = append(messages, bytes.Trim(message, "\x00"))
	}

	return helpers.ByteArraySeriesToString(messages), nil
}

func (s *hsmService) reconnect() (p11.Session, core.IError) {
	var err error
	var session p11.Session
	maxRetry := 4
	retry := 0

	for {
		session, err = helpers.NewHSMSession(s.ctx.ENV().Int(consts.ENVHSMSlot), s.ctx.ENV().String(consts.ENVHSMPin))
		retry = retry + 1
		if err == nil {
			break
		}
		if retry >= maxRetry {
			s.ctx.Log().Info(fmt.Sprintf("reconnecting to HSM failed after %v attepms", retry))
			return nil, s.ctx.NewError(emsgs.HSMSessionError(err), emsgs.HSMSessionError(err))
		}
		time.Sleep(200 * time.Millisecond)
	}

	s.ctx.Log().Info(fmt.Sprintf("reconnecting to HSM successfully after %v attepms", retry))

	return session, nil
}

func (s *hsmService) getPublicKey(session p11.Session) (*p11.PublicKey, error) {
	publicKeyTemplate := []*pkcs11.Attribute{pkcs11.NewAttribute(pkcs11.CKA_CLASS, pkcs11.CKO_PUBLIC_KEY)}
	pubilcKeyObject, err := session.FindObject(publicKeyTemplate)
	if err != nil {
		return nil, err
	}
	publicKey := p11.PublicKey(pubilcKeyObject)
	return &publicKey, nil
}

func (s *hsmService) getPrivateKey(session p11.Session) (*p11.PrivateKey, error) {
	privateKeyTemplate := []*pkcs11.Attribute{pkcs11.NewAttribute(pkcs11.CKA_CLASS, pkcs11.CKO_PRIVATE_KEY)}
	privateKeyObject, err := session.FindObject(privateKeyTemplate)
	if err != nil {
		return nil, err
	}
	privateKey := p11.PrivateKey(privateKeyObject)
	return &privateKey, nil
}

func (s *hsmService) encrypt(plaintext []byte) ([]byte, core.IError) {
	session, ok := s.ctx.GetData(consts.ContextKeyHSMSession).(p11.Session)
	if !ok {
		err := errors.New("cannot get session from context")
		return nil, s.ctx.NewError(emsgs.HSMSessionError(err), emsgs.HSMSessionError(err))
	}

	publicKey, err := s.getPublicKey(session)
	if err != nil {
		newSesion, ierr := s.reconnect()
		if ierr != nil {
			return nil, s.ctx.NewError(ierr, ierr)
		}
		session = newSesion
		s.ctx.SetData(consts.ContextKeyHSMSession, newSesion)
		publicKey, err = s.getPublicKey(newSesion)
		if err != nil {
			return nil, s.ctx.NewError(emsgs.HSMObjectError(err), emsgs.HSMObjectError(err))
		}
	}

	mechanism := pkcs11.NewMechanism(pkcs11.CKM_RSA_PKCS_OAEP, pkcs11.NewOAEPParams(pkcs11.CKM_SHA256, pkcs11.CKG_MGF1_SHA256, pkcs11.CKZ_DATA_SPECIFIED, make([]byte, 0)))

	cipher, err := publicKey.Encrypt(*mechanism, plaintext)
	if err != nil {
		return nil, s.ctx.NewError(emsgs.HSMRSACryptographyError(err), emsgs.HSMRSACryptographyError(err))
	}

	return cipher, nil
}

func (s *hsmService) decrypt(cipher []byte) ([]byte, core.IError) {
	session, ok := s.ctx.GetData(consts.ContextKeyHSMSession).(p11.Session)
	if !ok {
		err := errors.New("cannot get session from context")
		return nil, s.ctx.NewError(emsgs.HSMSessionError(err), emsgs.HSMSessionError(err))
	}

	privateKey, err := s.getPrivateKey(session)
	if err != nil {
		newSesion, ierr := s.reconnect()
		if ierr != nil {
			return nil, s.ctx.NewError(ierr, ierr)
		}
		session = newSesion
		s.ctx.SetData(consts.ContextKeyHSMSession, newSesion)
		privateKey, err = s.getPrivateKey(newSesion)
		if err != nil {
			return nil, s.ctx.NewError(emsgs.HSMObjectError(err), emsgs.HSMObjectError(err))
		}
	}

	mechanism := pkcs11.NewMechanism(pkcs11.CKM_RSA_PKCS_OAEP, pkcs11.NewOAEPParams(pkcs11.CKM_SHA256, pkcs11.CKG_MGF1_SHA256, pkcs11.CKZ_DATA_SPECIFIED, make([]byte, 0)))
	message, err := privateKey.Decrypt(*mechanism, cipher)
	if err != nil {
		return nil, s.ctx.NewError(emsgs.HSMRSACryptographyError(err), emsgs.HSMRSACryptographyError(err))
	}

	return message, nil
}
