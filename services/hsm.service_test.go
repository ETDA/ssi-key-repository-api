package services

import (
	"github.com/stretchr/testify/suite"
	core "ssi-gitlab.teda.th/ssi/core"
	"testing"
)

type HSMServiceTestSuite struct {
	suite.Suite
	ctx        *core.ContextMock
	hs         IHSMService
	privateKey string
}

func TestHSMServiceTestSuite(t *testing.T) {
	suite.Run(t, new(HSMServiceTestSuite))
}

func (h *HSMServiceTestSuite) SetupTest() {
}

func (h *HSMServiceTestSuite) TestHSMServiceTestSuite_Encrypt_ExpectSuccess() {
}

func (h *HSMServiceTestSuite) TestHSMServiceTestSuite_Encrypt_ExpectError() {
}
