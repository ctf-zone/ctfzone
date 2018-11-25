package crypto

import "golang.org/x/crypto/bcrypt"

// HashPassword hashes given passwords.
func HashPassword(password string) (string, error) {

	hash, err := bcrypt.GenerateFromPassword(
		[]byte(password),
		bcrypt.DefaultCost,
	)
	if err != nil {
		return "", err
	}

	return string(hash), nil
}

// CheckPassword compares given password and hash.
func CheckPassword(hash, password string) error {
	return bcrypt.CompareHashAndPassword(
		[]byte(hash),
		[]byte(password),
	)
}
