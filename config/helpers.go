package config

import (
	"encoding/hex"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

func (c Game) StartTime() time.Time {
	t, err := time.Parse(time.RFC3339, c.Start)
	if err != nil {
		panic(err)
	}
	return t
}

func (c Game) EndTime() time.Time {
	t, err := time.Parse(time.RFC3339, c.End)
	if err != nil {
		panic(err)
	}
	return t
}

func (c Session) SecretBytes() []byte {
	b, err := hex.DecodeString(c.Secret)
	if err != nil {
		panic(err)
	}
	return b
}

func (c Server) BaseURL() string {
	url := ""

	if c.TLS.Enabled {
		url += "https://"
	} else {
		url += "http://"
	}

	url += c.Domain

	if c.IncludePort {
		if c.TLS.Enabled {
			url += fmt.Sprintf(":%d", c.TLS.Port)
		} else {
			url += fmt.Sprintf(":%d", c.Port)
		}
	}

	return url
}

func (c Server) AdminBaseURL() string {
	url := ""

	if c.TLS.Enabled {
		url += "https://"
	} else {
		url += "http://"
	}

	url += fmt.Sprintf("admin.%s", c.Domain)

	if c.IncludePort {
		if c.TLS.Enabled {
			url += fmt.Sprintf(":%d", c.TLS.Port)
		} else {
			url += fmt.Sprintf(":%d", c.Port)
		}
	}

	return url
}

func (c Files) Dir() string {
	if filepath.IsAbs(c.Path) {
		return c.Path
	}

	workDir, _ := os.Getwd()

	return filepath.Join(workDir, c.Path)
}
