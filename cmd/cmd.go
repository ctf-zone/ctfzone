package cmd

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/alexedwards/scs"
	"github.com/alexedwards/scs/stores/cookiestore"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/ctf-zone/ctfzone/config"
	"github.com/ctf-zone/ctfzone/controllers"
	"github.com/ctf-zone/ctfzone/internal/mailer"
	"github.com/ctf-zone/ctfzone/models"
	"github.com/ctf-zone/ctfzone/models/migrations"
)

var (
	cfgFile string
	debug   bool
)

var Cmd = &cobra.Command{
	Use:   "ctfzone",
	Short: "CTFZone server starter",
	Long:  "CTFZone server starter",
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		if debug {
			log.SetFormatter(&log.TextFormatter{})
		} else {
			log.SetFormatter(&log.JSONFormatter{})
		}

		if err := viper.ReadInConfig(); err == nil {
			log.Infof("Using config file: %s", viper.ConfigFileUsed())
		}
	},
	Run: func(cmd *cobra.Command, args []string) {

		var (
			err error
			cfg config.Config
		)

		if err := viper.Unmarshal(&cfg); err != nil {
			log.Fatal(errors.Wrap(err, "config: fail to unmarshal"))
		}

		if err := cfg.Validate(); err != nil {
			log.Fatal(err)
		}

		var (
			m mailer.Sender
		)

		if cfg.Email.Enabled {
			m, err = mailer.New(cfg.Email.DSN, cfg.Email.Templates.Base, cfg.Email.Templates.Ext)
			if err != nil {
				log.Fatal(err)
			}
		}

		var (
			db *models.Repository
		)

		if db, err = models.New(cfg.DB.DSN); err != nil {
			log.Fatal(err)
		}
		defer db.Close()

		if debug {
			db.EnableLogging()
		}

		if err := migrations.Up(cfg.DB.DSN); err != nil {
			log.Fatal(err)
		}

		var (
			s scs.Store
		)

		s = cookiestore.New(cfg.Session.SecretBytes())

		router := controllers.Router(&cfg, db, m, s)

		log.Info(fmt.Sprintf("Starting server on %s:%d",
			cfg.Server.Addr, cfg.Server.Port))

		srv := http.Server{
			Handler: router,
		}

		if cfg.Server.TLS.Enabled {

			go func() {
				addr := fmt.Sprintf("%s:%d", cfg.Server.Addr, cfg.Server.Port)
				log.Fatal(http.ListenAndServe(addr, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					http.Redirect(w, r, cfg.Server.BaseURL()+r.RequestURI, http.StatusMovedPermanently)
				})))
			}()

			log.Info(fmt.Sprintf("Starting TLS server on %s:%d",
				cfg.Server.Addr, cfg.Server.TLS.Port))

			srv.Addr = fmt.Sprintf("%s:%d", cfg.Server.Addr, cfg.Server.TLS.Port)

			srv.TLSConfig = &tls.Config{
				MinVersion: tls.VersionTLS12,
				CurvePreferences: []tls.CurveID{
					tls.CurveP521,
					tls.CurveP384,
					tls.CurveP256,
				},
				PreferServerCipherSuites: true,
				CipherSuites: []uint16{
					tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
					tls.TLS_ECDHE_RSA_WITH_AES_256_CBC_SHA,
					tls.TLS_RSA_WITH_AES_256_GCM_SHA384,
					tls.TLS_RSA_WITH_AES_256_CBC_SHA,
				},
			}

			srv.TLSNextProto = make(map[string]func(*http.Server, *tls.Conn, http.Handler), 0)

			log.Fatal(srv.ListenAndServeTLS(cfg.Server.TLS.Cert, cfg.Server.TLS.Key))

		} else {
			srv.Addr = fmt.Sprintf("%s:%d", cfg.Server.Addr, cfg.Server.Port)

			log.Fatal(srv.ListenAndServe())
		}
	},
}

func Execute() {
	if err := Cmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	Cmd.Flags().BoolVarP(&debug, "debug", "d", false, "Debug mode")
	Cmd.Flags().StringVarP(&cfgFile, "config", "c", "", "Config file")

	Cmd.Flags().IntP("port", "p", 1337, "Port server listen on")

	viper.BindPFlag("server.port", Cmd.Flags().Lookup("port"))
}

// initConfig reads in config file.
func initConfig() {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	}

	viper.AutomaticEnv()
	viper.SetEnvPrefix("ctf")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	viper.SetConfigName("ctfzone")
	viper.AddConfigPath(".")
	viper.AddConfigPath("/etc")
}
