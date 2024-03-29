package router

import (
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/ohmpatel1997/findhotel-geolocation/integration/log"
	"net/http"
	"os"

	"github.com/go-chi/chi"
)

const (
	defaultCertLocation = "./ssl/cert.pem"
	defaultKeyLocation  = "./ssl/key.pem"

	defaultHealthCheckPath = "/healthcheck.html"
)

// Router interface, a subset of chi with some convenience methods
type Router interface {
	Delete(string, http.HandlerFunc, ...func(http.Handler) http.Handler)
	Get(string, http.HandlerFunc, ...func(http.Handler) http.Handler)
	Patch(string, http.HandlerFunc, ...func(http.Handler) http.Handler)
	Post(string, http.HandlerFunc, ...func(http.Handler) http.Handler)
	Put(string, http.HandlerFunc, ...func(http.Handler) http.Handler)
	Options(string, http.HandlerFunc, ...func(http.Handler) http.Handler)

	Route(string, func(r Router)) Router

	Handle(string, http.Handler)
	HandleFunc(string, http.HandlerFunc)
	With(middlewares ...func(http.Handler) http.Handler) Router

	ListenAndServeTLS(string, *tls.Config) error
	ServeHTTP(http.ResponseWriter, *http.Request)
}

type router struct {
	chi *chi.Mux
}

// NewBasicRouter is a basic router without authorization for back compat
func NewBasicRouter() Router {
	rchi := chi.NewRouter()
	rchi.Use(LoggerAndRecover)

	return router{
		chi: rchi,
	}
}

func (r router) With(middlewares ...func(http.Handler) http.Handler) Router {
	r.chi = r.chi.With(middlewares...).(*chi.Mux)
	return r
}

func (r router) Delete(p string, h http.HandlerFunc, middlewares ...func(http.Handler) http.Handler) {
	r.chi.With(middlewares...).Delete(p, h)
}

func (r router) Get(p string, h http.HandlerFunc, middlewares ...func(http.Handler) http.Handler) {
	r.chi.With(middlewares...).Get(p, h)
}

func (r router) Patch(p string, h http.HandlerFunc, middlewares ...func(http.Handler) http.Handler) {
	r.chi.With(middlewares...).Patch(p, h)
}

func (r router) Post(p string, h http.HandlerFunc, middlewares ...func(http.Handler) http.Handler) {
	r.chi.With(middlewares...).Post(p, h)
}

func (r router) Put(p string, h http.HandlerFunc, middlewares ...func(http.Handler) http.Handler) {
	r.chi.With(middlewares...).Put(p, h)
}

func (r router) Options(p string, h http.HandlerFunc, middlewares ...func(http.Handler) http.Handler) {
	r.chi.With(middlewares...).Options(p, h)
}

func (r router) Route(p string, fn func(r Router)) Router {
	nr := router{chi.NewRouter()} //get new router

	if fn != nil {
		fn(nr) //register the sub path
	}

	r.Mount(p, nr) //mount the new router
	return nr
}

func (r router) Mount(p string, h http.Handler) {
	r.chi.Mount(p, h)
}

func (r router) Handle(p string, h http.Handler) {
	r.chi.Handle(p, h)
}

func (r router) HandleFunc(p string, h http.HandlerFunc) {
	r.chi.HandleFunc(p, h)
}

func (r router) ListenAndServeTLS(listenPort string, config *tls.Config) error {
	if listenPort == "" {
		return errors.New("invalid or missing listen port")
	}

	server := &http.Server{
		Addr:      fmt.Sprintf(":%s", listenPort),
		Handler:   r.chi,
		TLSConfig: config,
	}

	if config != nil {
		server.TLSConfig.BuildNameToCertificate()
	}

	if _, err := os.Stat(defaultCertLocation); os.IsNotExist(err) {
		return server.ListenAndServe()
	}

	return server.ListenAndServe()
}

func (r router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	r.chi.ServeHTTP(w, req)
}

//Response is all the info we need to properly render json ResponseWriter, Data, Logger, Status
type Response struct {
	Writer http.ResponseWriter
	Data   interface{}
	Logger log.Logger
	Status int
}

type HttpError struct {
	Message string `json:"message"`
	Status  int32  `json:"status"`
}

func NewHttpError(message string, status int32) *HttpError {
	return &HttpError{message, status}
}

func RenderJSON(r Response) {
	var j []byte
	var err error

	j, err = json.Marshal(r.Data)

	r.Writer.Header().Set("Content-Type", "application/json")

	if err != nil {
		r.Logger.ErrorD("Error marshalling repsonse data", log.Fields{"err": err.Error()})

		r.Writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	if r.Status > 0 {
		r.Writer.WriteHeader(r.Status)
	}

	r.Writer.Write(j)
}

func GetLogger(r *http.Request) log.Logger {
	ctx := r.Context()

	l, ok := ctx.Value(loggerKey).(log.Logger)

	//If the logger isn't present in the request(this should never happen)
	//Let's add a logger to the request, but log a warning
	if !ok {
		l = log.NewLogger()

		*r = *addLoggerContextToRequest(l, r)

		l.WarnD("Logger missing from request context", log.Fields{"path": r.URL.Path})
	}

	return l
}
