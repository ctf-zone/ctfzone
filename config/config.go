package config

import (
	"time"

	"github.com/spf13/viper"
)

type CSRF struct {
	Lifetime time.Duration `json:"lifetime"`
}

type DB struct {
	DSN string `json:"dsn"`
}

type Email struct {
	Enabled   bool   `json:"enabled"`
	DSN       string `json:"dsn"`
	Templates struct {
		Base string `json:"base"`
		Ext  string `json:"ext"`
	} `json:"templates"`
}

type Scoring struct {
	Type    string `json:"type"`
	Dynamic struct {
		Min   int     `json:"min"`
		Max   int     `json:"max"`
		Coeff float64 `json:"coeff"`
	} `json:"dynamic"`
}

type Game struct {
	Start      string   `json:"start"`
	End        string   `json:"end"`
	Scoring    Scoring  `json:"scoring"`
	Categories []string `json:"categories"`
}

type ReCaptcha struct {
	Enabled bool   `json:"enabled"`
	Secret  string `json:"secret"`
	Header  string `json:"header"`
}

type Session struct {
	Secret   string        `json:"secret"`
	Lifetime time.Duration `json:"lifetime"`
}

type Files struct {
	Path    string `json:"path"`
	Prefix  string `json:"prefix"`
	MaxSize int64  `json:"maxSize"`
}

type TLS struct {
	Enabled bool   `json:"enabled"`
	Port    int    `json:"port"`
	Key     string `json:"key"`
	Cert    string `json:"cert"`
}

type Server struct {
	Addr        string `json:"addr"`
	Port        int    `json:"port"`
	Domain      string `json:"domain"`
	IncludePort bool   `json:"includePort"`

	TLS TLS `json:"tls"`
}

type Public struct {
	StaticPath string `json:"staticPath"`
}

type Admin struct {
	StaticPath string `json:"staticPath"`
	Password   string `json:"password"`
}

type Config struct {
	CSRF      CSRF      `json:"csrf"`
	DB        DB        `json:"db"`
	Email     Email     `json:"email"`
	Game      Game      `json:"game"`
	ReCaptcha ReCaptcha `json:"reCaptcha"`
	Session   Session   `json:"session"`
	Files     Files     `json:"files"`
	Server    Server    `json:"server"`
	Admin     Admin     `json:"admin"`
	Public    Public    `json:"public"`
}

func init() {

	// DB
	viper.SetDefault("db.dsn", "")

	// CSRF
	viper.SetDefault("csrf.lifetime", "2h")

	// Email
	viper.SetDefault("email.enabled", false)
	viper.SetDefault("email.dsn", "")
	viper.SetDefault("email.templates.base", "./templates/email/base.tmpl")
	viper.SetDefault("email.templates.ext", "./templates/email/ext/*.tmpl")

	// Game
	viper.SetDefault("game.start", "")
	viper.SetDefault("game.end", "")
	viper.SetDefault("game.scoring.type", "classic")
	viper.SetDefault("game.categories", []string{
		"crypto",
		"misc",
		"pwn",
		"recon",
		"reverse",
		"web",
		"forensics",
		"stegano",
	})
	viper.SetDefault("game.scoring.dynamic.min", 100)
	viper.SetDefault("game.scoring.dynamic.max", 500)
	viper.SetDefault("game.scoring.dynamic.coeff", 0.925)

	// ReCaptcha
	viper.SetDefault("recaptcha.enabled", false)
	viper.SetDefault("recaptcha.secret", "")
	viper.SetDefault("recaptcha.header", "X-G-Recaptcha-Response")

	// Files
	viper.SetDefault("files.path", "./files")
	viper.SetDefault("files.prefix", "/files")
	viper.SetDefault("files.maxSize", 32<<20)

	// Server
	viper.SetDefault("server.addr", "0.0.0.0")
	viper.SetDefault("server.port", 8080)
	viper.SetDefault("server.domain", "")
	viper.SetDefault("server.includePort", false)

	viper.SetDefault("server.tls.enabled", false)
	viper.SetDefault("server.tls.port", 8443)
	viper.SetDefault("server.tls.key", "")
	viper.SetDefault("server.tls.cert", "")

	// Session
	viper.SetDefault("session.secret", "")
	viper.SetDefault("session.lifetime", "24h")

	// Admin
	viper.SetDefault("admin.password", "")
	viper.SetDefault("admin.staticPath", "./static/admin")

	// Public
	viper.SetDefault("public.staticPath", "./static/public")
}
