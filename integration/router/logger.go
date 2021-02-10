package router

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/ohmpatel1997/findhotel-geolocation/integration/log"
)

const loggerKey = ctxKey("rlogger")

type ctxKey string

// Logger is a middleware used for all requests, see NewRouter
func LoggerAndRecover(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		l := log.NewLogger()

		r = addLoggerContextToRequest(l, r)

		sw := statusWriter{ResponseWriter: w}

		defer func(l log.Logger, r *http.Request) {
			err := recover()
			if err != nil {
				f := log.Fields{
					// err value from recover can be a non-error type
					"error":  fmt.Sprintf("%v", err),
					"host":   r.Host,
					"method": r.Method,
					"path":   r.URL.Path,
					"status": http.StatusInternalServerError,
				}

				l.ErrorD("ROUTER ERROR", f)

				jsonBody, _ := json.Marshal(map[string]string{
					"error": "There was an internal server error",
				})

				w.WriteHeader(http.StatusInternalServerError)
				w.Write(jsonBody)
			}
		}(l, r)

		start := time.Now()

		next.ServeHTTP(&sw, r)

		duration := time.Now().Sub(start)

		l.InfoD("ACCESS", log.Fields{
			"host":           r.Host,
			"method":         r.Method,
			"path":           r.URL.Path,
			"status":         sw.status,
			"content_length": sw.length,
			"duration":       duration,
		})
	}

	return http.HandlerFunc(fn)
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

func addLoggerContextToRequest(l log.Logger, r *http.Request) *http.Request {
	ctx := r.Context()
	ctx = context.WithValue(ctx, loggerKey, l)

	return r.WithContext(ctx)
}
