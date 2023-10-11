package server

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gorilla/mux"
	"go.uber.org/zap"
	"io"
	"net/http"
	"pismo-ledger-service/config"
	"pismo-ledger-service/errs"
	"pismo-ledger-service/pkg/log"
	"pismo-ledger-service/pkg/repository"
)

const (
	CT     = "Content-Type"
	CTJson = "application/json"
)

// Server struct
type Server struct {
	svc repository.DBOps
	cfg *config.Config
}

// NewServer constructor
func NewServer(cfg *config.Config, s repository.DBOps) *Server {
	return &Server{cfg: cfg, svc: s}
}

// Run server
func (h *Server) Run() *mux.Router {

	r := mux.NewRouter()

	r.HandleFunc("/accounts", h.createAccount).Methods(http.MethodPost)
	r.HandleFunc("/accounts/{account_id}", h.getAccount).Methods(http.MethodGet)
	r.HandleFunc("/transactions", h.createTransaction).Methods(http.MethodPost)

	return r
}

// GetParams unmarshalls request body to required struct
// Param - interface object in which the request body is to be unmarshalled
// Param - http request object from which the request is to be unmarshalled to desired object
func (h *Server) GetParams(o interface{}, request *http.Request) (err error) {
	ct := getContentType(request)
	if ct != CTJson {
		return errors.New("unsupported media type")
	}

	body, err := io.ReadAll(request.Body)
	if err != nil {
		return errors.New("error reading request body")
	}

	// Restore the io.ReadCloser to its original state
	request.Body = io.NopCloser(bytes.NewBuffer(body))

	if len(body) < 1 {
		return errs.ErrorEmptyBodyContent
	}

	err = json.Unmarshal(body, o)
	if err != nil {
		return errs.ErrorRequestBodyInvalid
	}

	return
}

// GetContentType ...
func getContentType(req *http.Request) (ct string) {
	ct = req.Header.Get(CT)
	return
}

// FormatException - formats the application exception and returns error response
func (h *Server) FormatException(r http.ResponseWriter, err error) {
	h.JSON(r, http.StatusBadRequest, errs.FormatErrorResponse(err))
}

// JSON sends a JSON response body
func (h *Server) JSON(r http.ResponseWriter, code int, content interface{}) {
	if fmt.Sprint(content) == "[]" {
		emptyResponse, _ := json.Marshal(make([]int64, 0))
		Output(r, code, CTJson, emptyResponse)
		return
	}

	var b bytes.Buffer
	enc := json.NewEncoder(&b)
	enc.SetEscapeHTML(false)
	enc.Encode(content)
	Output(r, code, CTJson, b.Bytes())
}

// Output sets a full HTTP output detail
func Output(r http.ResponseWriter, code int, ctype string, content []byte) {
	log.Info("Response ", zap.Any("Message", string(content)))
	r.Header().Set("Content-Type", ctype)
	r.WriteHeader(code)
	r.Write(content)
}
