package wsresponses

import (
	"github.com/Horqu/zkp-communicator-backend/cmd/internal"
)

func ResponseRegisterSuccess() internal.Response {
	return internal.Response{
		Command: internal.ResponseRegisterSuccess,
		Data:    "OK",
	}
}
