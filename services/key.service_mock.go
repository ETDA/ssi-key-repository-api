package services

import (
	"github.com/stretchr/testify/mock"
	"gitlab.finema.co/finema/etda/key-repository-api/models"
	core "ssi-gitlab.teda.th/ssi/core"
)

type MockKeyService struct {
	mock.Mock
}

func NewMockKeyService() *MockKeyService {
	return &MockKeyService{}
}

func (m *MockKeyService) Find(id string) (*models.Key, core.IError) {
	args := m.Called(id)
	return args.Get(0).(*models.Key), core.MockIError(args, 1)
}

func (m *MockKeyService) Store(payload *KeyStorePayload) (*models.Key, core.IError) {
	args := m.Called(payload)
	return args.Get(0).(*models.Key), core.MockIError(args, 1)
}

func (m *MockKeyService) Generate() (*models.Key, core.IError) {
	args := m.Called()
	return args.Get(0).(*models.Key), core.MockIError(args, 1)
}

func (m *MockKeyService) Sign(payload *KeySignPayload) (string, string, core.IError) {
	args := m.Called(payload)
	return args.String(0), args.String(1), core.MockIError(args, 2)
}
