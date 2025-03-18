package wsresponses

import (
	"github.com/Horqu/zkp-communicator-backend/cmd/internal"
)

func ResponseRegisterPage() internal.Response {
	return internal.Response{
		Command: internal.ResponseRegisterPage,
		Data:    "",
	}
}
