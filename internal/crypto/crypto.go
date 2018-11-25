package crypto

import "errors"

var (
	ErrInvalidHash = errors.New("Invalid hash format")
	ErrMismatch    = errors.New("Hash mismatch")
)
