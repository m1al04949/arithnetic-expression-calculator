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

// ComparePasswordWithHash сравнивает переданный пароль с хешем пароля
func ComparePasswordWithHash(password, hashedPassword string) bool {

	hashedPasswordBytes := []byte(hashedPassword)

	err := bcrypt.CompareHashAndPassword(hashedPasswordBytes, []byte(password))

	return err == nil
}
