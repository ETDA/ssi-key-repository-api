package emsgs

import (
	"net/http"

	core "ssi-gitlab.teda.th/ssi/core"
)

func HSMInitializeError(err error) core.IError {
	return &core.Error{
		Status:  http.StatusInternalServerError,
		Code:    "HSM_INITIALIZE_ERROR",
		Message: err.Error(),
	}
}

func HSMSlotError(err error) core.IError {
	return &core.Error{
		Status:  http.StatusInternalServerError,
		Code:    "HSM_SLOT_ERROR",
		Message: err.Error(),
	}
}

func HSMSessionError(err error) core.IError {
	return &core.Error{
		Status:  http.StatusInternalServerError,
		Code:    "HSM_SESSION_ERROR",
		Message: err.Error(),
	}
}

func HSMLoginError(err error) core.IError {
	return &core.Error{
		Status:  http.StatusInternalServerError,
		Code:    "HSM_LOGIN_ERROR",
		Message: err.Error(),
	}
}

func HSMLogoutError(err error) core.IError {
	return &core.Error{
		Status:  http.StatusInternalServerError,
		Code:    "HSM_LOGOUT_ERROR",
		Message: err.Error(),
	}
}

func HSMObjectError(err error) core.IError {
	return &core.Error{
		Status:  http.StatusInternalServerError,
		Code:    "HSM_OBJECT_ERROR",
		Message: err.Error(),
	}
}

func HSMRSACryptographyError(err error) core.IError {
	return &core.Error{
		Status:  http.StatusInternalServerError,
		Code:    "HSM_RSA_CRYPTOGRAPHY_ERROR",
		Message: err.Error(),
	}
}
