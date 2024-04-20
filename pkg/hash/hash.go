package hash

import "golang.org/x/crypto/bcrypt"

// HashPassword хэширует пароль с помощью алгоритма bcrypt
func HashPassword(password string) (string, error) {

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	return string(hashedPassword), nil
}
