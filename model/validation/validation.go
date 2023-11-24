package validation

import (
	"github.com/asaskevich/govalidator"
)

func ValidateStruct(payload interface{}) error {
	_, err := govalidator.ValidateStruct(payload)

	if err != nil {
		return err
	}

	return nil
}
