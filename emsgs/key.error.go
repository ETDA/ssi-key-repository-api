package emsgs

import (
	"net/http"

	core "ssi-gitlab.teda.th/ssi/core"
)

var (
	KeyNotFoundError = core.Error{
		Status:  http.StatusNotFound,
		Code:    "KEY_NOT_FOUND",
		Message: "key is not found",
	}

	GenerateKeyError = core.Error{
		Status:  http.StatusInternalServerError,
		Code:    "GENERATE_KEY_ERROR",
		Message: "generate key error",
	}

	UnsupportedSigningAlgorithm = core.Error{
		Status:  http.StatusBadRequest,
		Code:    "UNSUPPORTED_APGORITHM",
		Message: "Proived key is unsupported",
	}
)
