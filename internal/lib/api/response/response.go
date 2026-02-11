package response

import (
	"app/internal/storage"
	"strings"

	"github.com/go-playground/validator/v10"
)

type ValidError struct {
	ValidErr string
	Status   int
}

func (v *ValidError) Error() string {
	return v.ValidErr
}

func ValidationError(errs validator.ValidationErrors) *ValidError {
	var errorMsg []string

	for _, err := range errs {
		switch err.ActualTag() {
		case "uuid4":
			errorMsg = append(errorMsg, err.Field()+" must be a valid UUIDv4")
		case "required":
			errorMsg = append(errorMsg, err.Field()+" is required")
		case "email":
			errorMsg = append(errorMsg, err.Field()+" must be a valid email")
		case "min":
			errorMsg = append(errorMsg, err.Field()+" must be at least "+err.Param()+" characters long")
		case "max":
			errorMsg = append(errorMsg, err.Field()+" must be at most "+err.Param()+" characters long")
		default:
			errorMsg = append(errorMsg, err.Field()+" is invalid")
		}
	}

	validateErr := ValidError{
		ValidErr: strings.Join(errorMsg, ", "),
		Status:   storage.StatusUnprocessableEntity,
	}

	return &validateErr
}
