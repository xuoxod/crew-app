package config

import (
	"log"
	"text/template"

	"github.com/alexedwards/scs/v2"
)

// App Configuration
type AppConfig struct {
	UseCache      bool
	TemplateCache map[string]*template.Template
	InfoLog       *log.Logger
	ErrorLog      *log.Logger
	InProduction  bool
	Session       *scs.SessionManager
	DBConnection  string
}
