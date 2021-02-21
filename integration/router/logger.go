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

func addLoggerContextToRequest(l log.Logger, r *http.Request) *http.Request {
	ctx := r.Context()
	ctx = context.WithValue(ctx, loggerKey, l)

	return r.WithContext(ctx)
}
