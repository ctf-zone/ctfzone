package crypto

import (
	"crypto/sha256"
	"crypto/subtle"
	"encoding/hex"
	"fmt"
)

func CheckFlag(hash, flag string) error {
	expectedHash, err := hex.DecodeString(hash)
	if err != nil {
		return ErrInvalidHash
	}

	actualHash := sha256.Sum256([]byte(flag))

	if subtle.ConstantTimeCompare(
		expectedHash,
		actualHash[:],
	) != 1 {
		return ErrMismatch
	}

	return nil
}

func HashFlag(flag string) string {
	return fmt.Sprintf("%x", sha256.Sum256([]byte(flag)))
}
