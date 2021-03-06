package auth

import (
	"net/http"
	"strings"

	"go.uber.org/zap"
)

type application string

const (
	// OfficeApp indicates office.move.mil
	OfficeApp application = "OFFICE"
	// MyApp indicates my.move.mil
	MyApp application = "MY"
)

// IsOfficeApp returns true iff the request is for the office.move.mil host
func (s *Session) IsOfficeApp() bool {
	return s.ApplicationName == OfficeApp
}

// IsMyApp returns true iff the request is for the my.move.mil host
func (s *Session) IsMyApp() bool {
	return s.ApplicationName == MyApp
}

// DetectorMiddleware detects which application we are serving based on the hostname
func DetectorMiddleware(logger *zap.Logger, myHostname string, officeHostname string) func(next http.Handler) http.Handler {
	logger.Info("Creating host detector", zap.String("myHost", myHostname), zap.String("officeHost", officeHostname))
	return func(next http.Handler) http.Handler {
		mw := func(w http.ResponseWriter, r *http.Request) {
			session := SessionFromRequestContext(r)
			parts := strings.Split(r.Host, ":")
			var appName application
			if strings.EqualFold(parts[0], myHostname) {
				appName = MyApp
			} else if strings.EqualFold(parts[0], officeHostname) {
				appName = OfficeApp
			} else {
				logger.Error("Bad hostname", zap.String("hostname", r.Host))
				http.Error(w, http.StatusText(400), http.StatusBadRequest)
				return
			}
			session.ApplicationName = appName
			session.Hostname = strings.ToLower(parts[0])
			next.ServeHTTP(w, r)
			return
		}
		return http.HandlerFunc(mw)
	}
}
