package auth

import "golang.org/x/crypto/bcrypt"

// HashPassword uses bcrypt in order to create a secure password hash.
func HashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

// CheckPasswordHash compares the password provided by the user to the hashed password stored in the database.
// Returns nil on success, or an error on failure.
func CheckPasswordHash(password, hash string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
}
