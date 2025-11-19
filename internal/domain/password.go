package domain

import (
	"database/sql/driver"
	"errors"

	"golang.org/x/crypto/bcrypt"
)

// Password represents a password for user entity.
type Password struct {
	hash []byte
}

// Set genera el hash de la contraseña
func (p *Password) Set(plain string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(plain), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	p.hash = hash
	return nil
}

// Compare compara el texto plano con el hash
func (p *Password) Compare(plain string) error {
	return bcrypt.CompareHashAndPassword(p.hash, []byte(plain))
}

// Value que GORM usa automáticamente
func (p Password) Value() (driver.Value, error) {
	if len(p.hash) == 0 {
		return nil, nil
	}
	return string(p.hash), nil
}

// Scan también va dentro del struct Password
func (p *Password) Scan(value any) error {
	if value == nil {
		p.hash = nil
		return nil
	}

	switch v := value.(type) {
	case []byte:
		p.hash = v
	case string:
		p.hash = []byte(v)
	default:
		return errors.New("unsupported type for Password scan")
	}

	return nil
}

func (p Password) HasHash() bool {
	return len(p.hash) > 0
}
