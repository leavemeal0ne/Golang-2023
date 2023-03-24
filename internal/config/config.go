package config

import (
	"github.com/alexedwards/scs/v2"
	"github.com/leavemeal0ne/Golang-2023/internal/driver"
	"html/template"
)

type Config struct {
	TemplateCache map[string]*template.Template
	UseCache      bool
	BD            *driver.DB
	Session       *scs.SessionManager
}
