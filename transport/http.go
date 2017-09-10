package transport

import (
	"github.com/gorilla/mux"
	"bitbucket.org/code_horse/pegasus/transport/http_transport"
)

func NewHttpTransporter(router *mux.Router) ITransporter {

	if router == nil {
		router = mux.NewRouter()
	}

	return &http_transport.Transporter{
		Router: router,
	}
}
