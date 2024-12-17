package validator

import (
	"fmt"
	"invoices-api/internal/models"

	"github.com/go-playground/validator/v10"
)

type InvoiceValidator struct {
	validate *validator.Validate
}

type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

func NewInvoiceValidator() *InvoiceValidator {
	validate := validator.New()
	validate.RegisterValidation("validStatus", validStatusCheck)
	return &InvoiceValidator{
		validate: validate,
	}
}

func (v *InvoiceValidator) ValidateInvoice(invoice *models.Invoice) []ValidationError {
	var errors []ValidationError

	err := v.validate.Struct(invoice)
	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			var element ValidationError
			element.Field = err.Field()
			element.Message = v.getErrorMsg(err)
			errors = append(errors, element)
		}
	}

	return errors
}

func (v *InvoiceValidator) ValidatePartialUpdate(updates map[string]interface{}) []ValidationError {
	var errors []ValidationError

	for field, value := range updates {
		switch field {
		case "status":
			if status, ok := value.(string); ok {
				if status == "" {
					continue
				}
				if !isValidStatus(status) {
					errors = append(errors, ValidationError{
						Field:   "Status",
						Message: "Status must be one of: Paid, Pending, Unpaid",
					})
				}
			}
		case "amount":
			if amount, ok := value.(float64); ok {
				if amount <= 0 {
					errors = append(errors, ValidationError{
						Field:   "Amount",
						Message: "Amount must be greater than 0",
					})
				}
			}
		case "service_name":
			if name, ok := value.(string); ok {
				if name != "" && len(name) < 2 {
					errors = append(errors, ValidationError{
						Field:   "ServiceName",
						Message: "Service name must be at least 2 characters",
					})
				}
			}
		case "invoice_number":
			if num, ok := value.(float64); ok {
				if num <= 0 {
					errors = append(errors, ValidationError{
						Field:   "InvoiceNumber",
						Message: "Invoice number must be greater than 0",
					})
				}
			}
		}
	}

	return errors
}

func (v *InvoiceValidator) getErrorMsg(err validator.FieldError) string {
	switch err.Tag() {
	case "required":
		return fmt.Sprintf("%s is required", err.Field())
	case "min":
		return fmt.Sprintf("%s must be greater than %s", err.Field(), err.Param())
	case "validStatus":
		return fmt.Sprintf("%s must be one of: Paid, Pending, Unpaid", err.Field())
	default:
		return fmt.Sprintf("%s is not valid", err.Field())
	}
}

func validStatusCheck(fl validator.FieldLevel) bool {
	status := fl.Field().String()
	return isValidStatus(status)
}

func isValidStatus(status string) bool {
	validStatuses := map[string]bool{
		"Paid":    true,
		"Pending": true,
		"Unpaid":  true,
	}
	return validStatuses[status]
}
