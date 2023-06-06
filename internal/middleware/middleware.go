package middleware

import (
	"github.com/alexedwards/scs/v2"
	"net/http"
)

var session *scs.SessionManager

func SetSession(s *scs.SessionManager) {
	session = s
}

// SessionLoad load and saves the sessions on every request
func SessionLoad(next http.Handler) http.Handler {
	return session.LoadAndSave(next)
}
