package validator

import (
	"strings"
	"time"

	"github.com/go-playground/validator"
	"golang.org/x/xerrors"
)

func Validation(target interface{}) error {
	validate := validator.New()
	if err := validate.RegisterValidation("date_format_check", DateFormatCheck); err != nil {
		return err
	}

	if err := validate.Struct(target); err != nil {
		var errorMessages []string
		for _, err := range err.(validator.ValidationErrors) {
			var errorMessage string
			fieldName := err.Field()

			switch fieldName {
			case "OpeningTime":
				typ := err.Tag()
				switch typ {
				case "required":
					errorMessage = "OpeningTimeは必須項目です"
				case "date_format_check":
					errorMessage = "OpeningTimeの日付形式が不正です"
				}
			case "ClosingTime":
				typ := err.Tag()
				switch typ {
				case "required":
					errorMessage = "OpeningTimeは必須項目です"
				case "date_format_check":
					errorMessage = "ClosingTimeの日付形式が不正です"
				}
			}
			errorMessages = append(errorMessages, errorMessage)
		}
		return xerrors.New(strings.Join(errorMessages, ","))
	}
	return nil
}

func DateFormatCheck(fl validator.FieldLevel) bool {
	if _, err := time.Parse("2006-01-02 15:04:05", fl.Field().String()); err == nil {
		return true
	}
	return false
}
