package auth

import "golang.org/x/crypto/bcrypt"

type BCrypt struct {
	cost int
}

func NewBCrypt(cost int) *BCrypt {
	return &BCrypt{cost: cost}
}

func (b *BCrypt) Hash(password string) (string, error) {
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), b.cost)
	if err != nil {
		return "", err
	}
	return string(hashed), nil
}

func (b *BCrypt) Compare(password, hash string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
}
