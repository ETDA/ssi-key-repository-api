// +build e2e

package services

import (
	"errors"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"gitlab.finema.co/finema/etda/key-repository-api/consts"
	"gitlab.finema.co/finema/etda/key-repository-api/emsgs"
	"gitlab.finema.co/finema/etda/key-repository-api/models"
	core "ssi-gitlab.teda.th/ssi/core"
	"ssi-gitlab.teda.th/ssi/core/errmsgs"
	"ssi-gitlab.teda.th/ssi/core/utils"
	"gorm.io/gorm"
	"regexp"
	"testing"
	"time"
)

type MockKeyData struct {
	ID         string
	PublicKey  string
	PrivateKey string
	Type       string
	CreatedAt  *time.Time
	UpdatedAt  *time.Time
}

type MockSignData struct {
	Message string
}

func NewMockKeyData() *MockKeyData {
	kp, _ := utils.GenerateKeyPair()
	return &MockKeyData{
		ID:         utils.GetUUID(),
		PublicKey:  kp.PublicKeyPem,
		PrivateKey: kp.PrivateKeyPem,
		Type:       string(consts.KeyTypeECDSA),
		CreatedAt:  utils.GetCurrentDateTime(),
		UpdatedAt:  utils.GetCurrentDateTime(),
	}
}

func NewMockSignData() *MockSignData {
	return &MockSignData{
		Message: "eyJvcGVyYXRpb24iOiJESURfUkVHSVNURVIiLCJwdWJsaWNfa2V5IjoiLS0tLS1CRUdJTiBQVUJMSUMgS0VZLS0tLS1cbk1Ga3dFd1lIS29aSXpqMENBUVlJS29aSXpqMERBUWNEUWdBRXo4c0c0Uk5iWlRMT2ZiV1NJelZSZHhITW9ZRExcbmhBZkpQWEdzZzl4UEtSamxvdE5LMDcxM3BkK1dhdXpBc0tUWExOaHNsSGIxWVZRWkwvK1FzWVZFekE9PVxuLS0tLS1FTkQgUFVCTElDIEtFWS0tLS0tXG4iLCJrZXlfdHlwZSI6IlNlY3AyNTZyMVZlcmlmaWNhdGlvbktleTIwMTgifQ==",
	}
}

type KeyServiceTestSuite struct {
	suite.Suite
	rCtx core.IContext
	mCtx *core.ContextMock
	rks  IKeyService
	rhs  IHSMService
	mks  *MockKeyService
	mhs  *MockHSMService
}

func TestKeyServiceTestSuite(t *testing.T) {
	suite.Run(t, new(KeyServiceTestSuite))
}

func (k *KeyServiceTestSuite) SetupSuite() {
	env := core.NewENVPath("./..")
	mysql, _ := core.NewDatabase(env.Config()).Connect()
	k.rCtx = core.NewContext(&core.ContextOptions{
		DB:  mysql,
		ENV: env,
	})
}

func (k *KeyServiceTestSuite) SetupTest() {
	k.mCtx = core.NewMockContext()
	k.rhs = NewHSMService(k.rCtx)
	k.rks = NewKeyService(k.rCtx, k.rhs)
	k.mhs = NewMockHSMService()
	k.mks = NewMockKeyService()

	k.mCtx.On("DB").Return(k.mCtx.MockDB.Gorm)
}

func (k *KeyServiceTestSuite) TestKeyService_Find_ExpectSuccess() {
	mockKeyData := NewMockKeyData()

	err := k.rCtx.DB().Create(models.Key{
		ID:                  mockKeyData.ID,
		PublicKey:           mockKeyData.PublicKey,
		PrivateKeyEncrypted: "encrypted_private_key",
		Type:                string(consts.KeyTypeECDSA),
		CreatedAt:           utils.GetCurrentDateTime(),
		UpdatedAt:           utils.GetCurrentDateTime(),
	}).Error
	k.NoError(err)

	key, ierr := k.rks.Find(mockKeyData.ID)
	k.NoError(ierr)
	k.NotNil(key)

	k.Equal(mockKeyData.ID, key.ID)
	k.Equal(mockKeyData.PublicKey, key.PublicKey)

	err = k.rCtx.DB().Delete(models.Key{}, "id = ?", mockKeyData.ID).Error
	k.NoError(err)
}

func (k *KeyServiceTestSuite) TestKeyService_Find_ExpectError() {
	mockKeyData := NewMockKeyData()

	k.rhs = NewHSMService(k.mCtx)
	k.rks = NewKeyService(k.mCtx, k.rhs)

	k.mCtx.MockDB.Mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `keys` WHERE id = ? ORDER BY `keys`.`id` LIMIT 1")).
		WithArgs(mockKeyData.ID).WillReturnError(gorm.ErrRecordNotFound)
	k.mCtx.On("NewError", mock.Anything, mock.Anything, mock.Anything).Return(emsgs.KeyNotFoundError).Once()

	key, ierr := k.rks.Find(mockKeyData.ID)
	k.Error(ierr)
	k.True(errors.Is(emsgs.KeyNotFoundError, ierr))
	k.Nil(key)

	k.mCtx.MockDB.Mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `keys` WHERE id = ? ORDER BY `keys`.`id` LIMIT 1")).
		WithArgs(mockKeyData.ID).WillReturnError(gorm.ErrUnsupportedDriver)
	k.mCtx.On("NewError", mock.Anything, mock.Anything, mock.Anything).Return(errmsgs.DBError).Once()

	key, ierr = k.rks.Find(mockKeyData.ID)
	k.Error(ierr)
	k.True(errors.Is(errmsgs.DBError, ierr))
	k.Nil(key)
}

func (k *KeyServiceTestSuite) TestKeyService_Store_ExpectSuccess() {
	mockKeyData := NewMockKeyData()

	key, ierr := k.rks.Store(&KeyStorePayload{
		PublicKey:  mockKeyData.PublicKey,
		PrivateKey: mockKeyData.PrivateKey,
		KeyType:    mockKeyData.Type,
	})
	k.NoError(ierr)
	k.NotNil(key)

	expectPrivateKey, ierr := k.rhs.Decrypt(key.PrivateKeyEncrypted)
	k.NoError(ierr)
	k.NotEmpty(expectPrivateKey)

	k.Equal(mockKeyData.PublicKey, key.PublicKey)
	k.NotEqual(mockKeyData.PrivateKey, key.PrivateKeyEncrypted)
	k.Equal(expectPrivateKey, mockKeyData.PrivateKey)
}

func (k *KeyServiceTestSuite) TestKeyService_Store_ExpectError() {
	mockKeyData := NewMockKeyData()

	k.rhs = NewHSMService(k.mCtx)
	k.rks = NewKeyService(k.mCtx, k.rhs)

	k.mCtx.MockDB.Mock.ExpectExec(regexp.QuoteMeta("INSERT INTO `keys` (`id`,`public_key`,`private_key_encrypted`,`type`,`created_at`,`updated_at`,`deleted_at`) VALUES (?,?,?,?,?,?,?)")).
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnError(gorm.ErrInvalidData)
	k.mCtx.On("NewError", mock.Anything, mock.Anything, mock.Anything).Return(errmsgs.DBError).Once()

	// Expect DBError
	key, ierr := k.rks.Store(&KeyStorePayload{
		PublicKey:  mockKeyData.PublicKey,
		PrivateKey: mockKeyData.PrivateKey,
	})
	k.Error(ierr)
	k.Equal(errmsgs.DBError.GetCode(), ierr.GetCode())
	k.Nil(key)

	// Expect DBError
	k.mhs.On("Encrypt", mockKeyData.PrivateKey).Return("", errmsgs.InternalServerError)
	k.rks = NewKeyService(k.mCtx, k.mhs)

	k.mCtx.MockDB.Mock.ExpectExec(regexp.QuoteMeta("INSERT INTO `keys` (`id`,`public_key`,`private_key_encrypted`,`type`,`created_at`,`updated_at`,`deleted_at`) VALUES (?,?,?,?,?,?,?)")).
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnError(gorm.ErrInvalidData)
	k.mCtx.On("NewError", mock.Anything, mock.Anything, mock.Anything).Return(errmsgs.InternalServerError).Once()

	key, ierr = k.rks.Store(&KeyStorePayload{
		PublicKey:  mockKeyData.PublicKey,
		PrivateKey: mockKeyData.PrivateKey,
	})
	k.Error(ierr)
	k.Equal(errmsgs.InternalServerError.GetCode(), ierr.GetCode())
	k.Nil(key)
}

func (k *KeyServiceTestSuite) TestKeyService_Generate_ExpectSuccess() {
	key, ierr := k.rks.Generate()
	k.NoError(ierr)
	k.NotNil(key)

	k.NotNil(key.ID)
	k.NotNil(key.PublicKey)
	k.NotNil(key.PrivateKeyEncrypted)
	k.NotNil(key.Type)
	k.NotNil(key.CreatedAt)
	k.NotNil(key.UpdatedAt)
	k.Nil(key.DeletedAt)
}

func (k *KeyServiceTestSuite) TestKeyService_Sign_ExpectSuccess() {
	mockKeyData := NewMockKeyData()
	mockSignData := NewMockSignData()

	encryptedPrivateKey, ierr := k.rhs.Encrypt(mockKeyData.PrivateKey)
	k.NoError(ierr)

	err := k.rCtx.DB().Create(models.Key{
		ID:                  mockKeyData.ID,
		PublicKey:           mockKeyData.PublicKey,
		PrivateKeyEncrypted: encryptedPrivateKey,
		Type:                string(consts.KeyTypeECDSA),
		CreatedAt:           utils.GetCurrentDateTime(),
		UpdatedAt:           utils.GetCurrentDateTime(),
	}).Error
	k.NoError(err)

	signature, ierr := k.rks.Sign(mockKeyData.ID, mockSignData.Message)
	k.NoError(ierr)
	k.NotEmpty(signature)

	valid, err := utils.VerifySignature(mockKeyData.PublicKey, signature, mockSignData.Message)

	k.NoError(err)
	k.True(valid)
	k.NotEmpty(signature)

	err = k.rCtx.DB().Delete(models.Key{}, "id = ?", mockKeyData.ID).Error
	k.NoError(err)
}

func (k *KeyServiceTestSuite) TestKeyService_Sign_ExpectError() {
	mockKeyData := NewMockKeyData()
	mockSignData := NewMockSignData()

	k.rhs = NewHSMService(k.mCtx)
	k.rks = NewKeyService(k.mCtx, k.rhs)

	k.mCtx.MockDB.Mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `keys` WHERE id = ? ORDER BY `keys`.`id` LIMIT 1")).
		WithArgs("invalid-ref-id").
		WillReturnError(gorm.ErrRecordNotFound)
	k.mCtx.On("NewError", mock.Anything, mock.Anything, mock.Anything).Return(emsgs.KeyNotFoundError).Twice()

	// Expect KeyNotFoundError at Find function
	signature, ierr := k.rks.Sign("invalid-ref-id", mockSignData.Message)
	k.Error(ierr)
	k.Equal(emsgs.KeyNotFoundError.GetCode(), ierr.GetCode())
	k.Empty(signature)

	// Expect InternalServerError at Encrypt function
	k.rhs = NewHSMService(k.rCtx)
	k.rks = NewKeyService(k.rCtx, k.rhs)

	encryptedPrivateKey, ierr := k.rhs.Encrypt(mockKeyData.PrivateKey)
	k.NoError(ierr)

	err := k.rCtx.DB().Create(models.Key{
		ID:                  mockKeyData.ID,
		PublicKey:           mockKeyData.PublicKey,
		PrivateKeyEncrypted: encryptedPrivateKey,
		Type:                string(consts.KeyTypeECDSA),
		CreatedAt:           utils.GetCurrentDateTime(),
		UpdatedAt:           utils.GetCurrentDateTime(),
	}).Error
	k.NoError(err)

	k.mhs.On("Decrypt", encryptedPrivateKey).Return("", errmsgs.InternalServerError)
	k.rks = NewKeyService(k.mCtx, k.mhs)

	k.mCtx.On("NewError", mock.Anything, mock.Anything, mock.Anything).Return(errmsgs.InternalServerError)
	signature, ierr = k.rks.Sign(mockKeyData.ID, mockSignData.Message)
	k.Error(ierr)
	k.Equal(errmsgs.InternalServerError.GetCode(), ierr.GetCode())
	k.Empty(signature)

	err = k.rCtx.DB().Delete(models.Key{}, "id = ?", mockKeyData.ID).Error
	k.NoError(err)
}
