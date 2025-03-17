package wsresponses

import (
	"github.com/Horqu/zkp-communicator-backend/cmd/internal"
)

func ResponseLoginPage() internal.Response {
	return internal.Response{
		Command: internal.ResponseLoginPage,
		Data:    "",
	}
}
