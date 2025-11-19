// Package validator provides functions to validate user input.
package validator

import (
	"fmt"
	"strings"
	"unicode"
)

func ValidatePassword(pwd string) error {
	var errs []string

	if len(pwd) < 8 {
		errs = append(errs, "muy corta")
	}

	var hasUpper, hasLower, hasNumber, hasSpecial bool
	for _, c := range pwd {
		switch {
		case unicode.IsUpper(c):
			hasUpper = true
		case unicode.IsLower(c):
			hasLower = true
		case unicode.IsNumber(c):
			hasNumber = true
		case unicode.IsPunct(c) || unicode.IsSymbol(c):
			hasSpecial = true
		}
	}

	if !hasUpper {
		errs = append(errs, "falta mayúscula")
	}
	if !hasLower {
		errs = append(errs, "falta minúscula")
	}
	if !hasNumber {
		errs = append(errs, "falta número")
	}
	if !hasSpecial {
		errs = append(errs, "falta símbolo")
	}

	if len(errs) > 0 {
		return fmt.Errorf("la contraseña no cumple los requisitos: %s", strings.Join(errs, ", "))
	}

	return nil
}
