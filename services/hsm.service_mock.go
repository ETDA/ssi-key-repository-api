package services

import (
	"github.com/stretchr/testify/mock"
	core "ssi-gitlab.teda.th/ssi/core"
)

type MockHSMService struct {
	mock.Mock
}

func NewMockHSMService() *MockHSMService {
	return &MockHSMService{}
}

func (m *MockHSMService) Decrypt(encryptedPrivateKey string) (string, core.IError) {
	args := m.Called(encryptedPrivateKey)
	return args.String(0), core.MockIError(args, 1)
}

func (m *MockHSMService) Encrypt(privateKey string) (string, core.IError) {
	args := m.Called(privateKey)
	return args.String(0), core.MockIError(args, 1)
}
