package mailer

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"html/template"
	"net/url"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/go-gomail/gomail"
	"github.com/pkg/errors"
)

const baseTemplateName = "email"

type Sender interface {
	Send(template, to string, data interface{}) error
}

type emailTemplate struct {
	*template.Template `json:"-"`
	From               string `json:"from"`
	Subject            string `json:"subject"`
}

type mailer struct {
	dialer    *gomail.Dialer
	templates map[string]*emailTemplate
}

func loadTemplates(basePath, extPath string) (map[string]*emailTemplate, error) {
	templates := make(map[string]*emailTemplate)

	// Base template.
	base := []string{basePath}

	// Find all extension templates.
	ext, err := filepath.Glob(extPath)
	if err != nil {
		return nil, err
	}

	for _, f := range ext {

		// Base and extension templates.
		files := append(base, f)

		name := strings.TrimSuffix(filepath.Base(f), filepath.Ext(f))
		tpl := &emailTemplate{template.Must(template.ParseFiles(files...)), "", ""}

		// Get "From" and "Subject" from meta section of extension template.
		if meta := tpl.Lookup("meta"); meta != nil {
			var b bytes.Buffer

			if err := meta.Execute(&b, nil); err != nil {
				return nil, errors.Wrap(err, fmt.Sprintf("fail to process meta-section of '%s'", f))
			}

			if err := json.Unmarshal(b.Bytes(), tpl); err != nil {
				return nil, errors.Wrap(err, fmt.Sprintf("fail to decode json from meta-section of '%s'", f))
			}
		}

		templates[name] = tpl
	}

	return templates, nil
}

func New(dsn string, basePath, extPath string) (Sender, error) {
	u, err := url.Parse(dsn)
	if err != nil {
		return nil, errors.Wrap(err, "mailer: fail to parse DSN")
	}

	if u.User == nil {
		return nil, errors.Wrap(err, "mailer: no username and password in DSN")
	}

	port, err := strconv.Atoi(u.Port())
	if err != nil {
		return nil, errors.Wrap(err, "mailer: fail to parse port")
	}

	password, _ := u.User.Password()

	d := gomail.NewDialer(u.Hostname(), port, u.User.Username(), password)
	d.TLSConfig = &tls.Config{InsecureSkipVerify: true}

	s, err := d.Dial()
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("mailer: fail to connect %s:%d", u.Hostname(), port))
	}

	_ = s.Close()

	t, err := loadTemplates(basePath, extPath)
	if err != nil {
		return nil, errors.Wrap(err, "mailer: fail to load emails templates")
	}

	return &mailer{d, t}, nil
}

func (m *mailer) Send(template, to string, data interface{}) error {

	tpl, ok := m.templates[template]
	if !ok {
		return fmt.Errorf("mailer: template '%s' not found", template)
	}

	msg := gomail.NewMessage()

	msg.SetHeader("Subject", tpl.Subject)
	msg.SetHeader("From", tpl.From)
	msg.SetHeader("To", to)

	var b bytes.Buffer
	if err := tpl.ExecuteTemplate(&b, baseTemplateName, data); err != nil {
		return errors.Wrap(err, "mailer: fail to process template")
	}

	msg.SetBody("text/html", string(b.Bytes()))

	s, err := m.dialer.Dial()
	if err != nil {
		return errors.Wrap(err, "mailer: fail connect to server")
	}

	if err := gomail.Send(s, msg); err != nil {
		return errors.Wrap(err, "mailer: fail to send email")
	}

	_ = s.Close()

	return nil
}
