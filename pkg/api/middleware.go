package api

import (
	"net/http"

	"github.com/sirupsen/logrus"
)

func loggingMiddlware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		fields := logrus.Fields{
			"method":         r.Method,
			"path":           r.URL.Path,
			"query":          r.URL.RawQuery,
			"origin":         r.Header.Get("Origin"),
			"referrer":       r.Header.Get("Referrer"),
			"user_agent":     r.Header.Get("User-Agent"),
			"content_length": r.ContentLength,
		}

		logrus.WithFields(fields).Info()

		next.ServeHTTP(rw, r)
	})
}
