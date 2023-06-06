package config

import (
	"github.com/alexedwards/scs/v2"
	"html/template"
	"net/http"
	"time"
)

type Config struct {
	TemplateCache map[string]*template.Template
	UseCache      bool
	Session       *scs.SessionManager
}

func NewAppConfig() *Config {
	app := Config{}
	session := scs.New()
	session.Lifetime = 24 * time.Hour
	session.Cookie.Persist = true
	session.Cookie.SameSite = http.SameSiteLaxMode
	session.Cookie.Secure = true
	app.Session = session
	return &app
}
