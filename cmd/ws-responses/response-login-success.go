package wsresponses

import (
	"github.com/Horqu/zkp-communicator-backend/cmd/internal"
)

func ResponseLoginSuccess(publicKey string) internal.Response {
	return internal.Response{
		Command: internal.ResponseLoginSuccess,
		Data:    "OK",
	}
}
