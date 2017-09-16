package transport

import (
	"bitbucket.org/code_horse/pegasus/transport/http_transport"
	"github.com/gorilla/mux"
	"net/http"
)

func NewHttpTransporter(router *mux.Router) ITransporter {

	if router == nil {
		router = mux.NewRouter()
	}

	return &http_transport.Transporter{
		Router: router,
	}
}

func StartHTTP(path string, Router http.Handler) {
	http.ListenAndServe("0.0.0.0:8900", Router)
}
