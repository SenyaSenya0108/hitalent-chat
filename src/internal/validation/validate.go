package validation

import (
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/go-playground/validator/v10"
)

type Service struct {
	validate *validator.Validate
}

var validate *validator.Validate

func InitValidator() {
	validate = validator.New()
}

func GetValidator() *Service {
	validationService := &Service{validate: validate}
	return validationService
}

func (v *Service) validationErrorToText(errors *string, e validator.FieldError) string {
	var textErr string

	switch e.Tag() {
	case "required":
		textErr = fmt.Sprintf("%s обязательное поле", e.Field())
	case "min":
		textErr = fmt.Sprintf("%s должен быть не менее %s символов", e.Field(), e.Param())
	case "max":
		textErr = fmt.Sprintf("%s должен быть не более %s символов", e.Field(), e.Param())
	case "email":
		textErr = fmt.Sprintf("%s должен быть валидным email-адресом", e.Field())
	case "gte":
		textErr = fmt.Sprintf("%s должен быть больше или равен %s", e.Field(), e.Param())
	case "lte":
		textErr = fmt.Sprintf("%s должен быть меньше или равен %s", e.Field(), e.Param())
	case "oneof":
		textErr = fmt.Sprintf("Невалидное значение %s у поля %s. Допустимые значения %s", e.Value(), e.Field(), e.Param())
	default:
		textErr = "Unknown error"
	}

	*errors += textErr + "\n"

	return *errors
}

func (v *Service) ValidationHttpRequest(s interface{}) (error, int) {
	err := v.validate.Struct(s)
	if err != nil {
		//TODO
		log.Println(err.Error())
		var validationErrors validator.ValidationErrors
		var textErrors string
		if errors.As(err, &validationErrors) {
			for _, err := range validationErrors {
				textErrors = v.validationErrorToText(&textErrors, err)
			}

			return fmt.Errorf("%v", textErrors), http.StatusBadRequest
		}

		var invalidValidationError *validator.InvalidValidationError
		if errors.As(err, &invalidValidationError) {
			textErrors = "Invalid request body"
			return fmt.Errorf("%v", textErrors), http.StatusInternalServerError
		}
	}

	return nil, 0
}
