package transport

import (
	"github.com/gorilla/mux"
	"net/http"
	"bitbucket.org/code_horse/pegasus/transport/tranhttp"
)

func NewHttpTransporter(router *mux.Router) ITransporter {

	if router == nil {
		router = mux.NewRouter()
	}

	return &tranhttp.Transporter{
		Router: router,
	}
}

func StartHTTP(path string, Router http.Handler) {
	http.ListenAndServe("0.0.0.0:8900", Router)
}
