package response

import (
	"strings"

	"github.com/go-playground/validator/v10"
)

func ValidationError(errs validator.ValidationErrors) string {
	var errorMsg []string

	for _, err := range errs {
		switch err.ActualTag() {
		case "required":
			errorMsg = append(errorMsg, err.Field()+" is required")
		case "email":
			errorMsg = append(errorMsg, err.Field()+" must be a valid email")
		case "min":
			errorMsg = append(errorMsg, err.Field()+" must be at least "+err.Param()+" characters long")
		default:
			errorMsg = append(errorMsg, err.Field()+" is invalid")
		}
	}
	return strings.Join(errorMsg, ", ")
}
