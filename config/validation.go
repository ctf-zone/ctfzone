package config

import (
	"errors"
	"os"
	"regexp"
	"time"

	"github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
)

func validateRFC3339(value interface{}) error {
	timeStr, _ := value.(string)

	if _, err := time.Parse(time.RFC3339, timeStr); err != nil {
		return err
	}

	return nil
}

func validateBcryptHash(value interface{}) error {
	validBcrypt := regexp.MustCompile(`^\$2a\$10\$.{53}$`)

	hash, _ := value.(string)

	if !validBcrypt.MatchString(hash) {
		return errors.New(`Invalid bcrypt hash`)
	}

	return nil
}

func validateFile(value interface{}) error {
	fpath, _ := value.(string)

	if _, err := os.Stat(fpath); os.IsNotExist(err) {
		return err
	}

	return nil
}

func (c DB) Validate() error {
	return validation.ValidateStruct(&c,
		validation.Field(&c.DSN, validation.Required, is.URL),
	)
}

func (c Email) Validate() error {
	if c.Enabled {
		return validation.ValidateStruct(&c,
			validation.Field(&c.DSN, validation.Required, is.URL),
		)
	}

	return nil
}

func (c Session) Validate() error {
	return validation.ValidateStruct(&c,
		validation.Field(&c.Secret, validation.Required, is.Hexadecimal, validation.Length(64, 64)))
}

func (c ReCaptcha) Validate() error {
	if c.Enabled {
		return validation.ValidateStruct(&c,
			validation.Field(&c.Secret, validation.Required),
		)
	}
	return nil
}

func (c Game) Validate() error {
	return validation.ValidateStruct(&c,
		validation.Field(&c.End, validation.Required, validation.By(validateRFC3339)),
	)
}

func (c TLS) Validate() error {
	if !c.Enabled {
		return nil
	}

	return validation.ValidateStruct(&c,
		validation.Field(&c.Key, validation.Required, validation.By(validateFile)),
		validation.Field(&c.Cert, validation.Required, validation.By(validateFile)),
	)
}

func (c Admin) Validate() error {
	return validation.ValidateStruct(&c,
		validation.Field(&c.Password, validation.Required, validation.By(validateBcryptHash)),
	)
}

func (c Server) Validate() error {
	return validation.ValidateStruct(&c,
		validation.Field(&c.Addr, validation.Required, is.IP),
		validation.Field(&c.Port, validation.Required, validation.Min(1), validation.Max(65535)),
		validation.Field(&c.Domain, validation.Required, is.Domain),
		validation.Field(&c.TLS),
	)
}

func (c Config) Validate() error {
	fieldsRules := []*validation.FieldRules{
		validation.Field(&c.Admin),
		validation.Field(&c.DB),
		validation.Field(&c.Email),
		validation.Field(&c.Game),
		validation.Field(&c.ReCaptcha),
		validation.Field(&c.Server),
		validation.Field(&c.Session),
	}

	return validation.ValidateStruct(&c, fieldsRules...)
}
