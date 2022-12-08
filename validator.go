package goracoon

import (
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/asaskevich/govalidator"
)

type Validation struct {
	Data   url.Values
	Errors map[string]string
}

// Validator initiates the validation
func (gr *Goracoon) Validator(data url.Values) *Validation {
	return &Validation{
		Data:   data,
		Errors: make(map[string]string),
	}
}

// Valid returns true if there are no erros in the validation
func (v *Validation) Valid() bool {
	return len(v.Errors) == 0
}

// Check used to make validations. Takes in a comparison statement, message key
// and message and if invalid adds the error to the validation
func (v *Validation) Check(ok bool, key string, message string) {
	if !ok {
		v.AddError(key, message)
	}
}

// AddError ads an error to the validation
func (v *Validation) AddError(key, message string) {
	if _, exists := v.Errors[key]; !exists {
		v.Errors[key] = message
	}
}

// HasField checks if a field exists in a HTTP request
func (v *Validation) HasField(field string, r *http.Request) bool {
	if x := r.Form.Get(field); x == "" {
		return false
	}
	return true
}

// HasRequiredFields checks if given fields exists in the HTTP request
func (v *Validation) HasRequiredFields(r *http.Request, fields ...string) {
	for _, field := range fields {
		value := r.Form.Get(field)
		if strings.TrimSpace(value) == "" {
			v.AddError(field, "this field cannot be blank")
		}
	}
}

// IsEmail validates if the given value is a valid email address
func (v *Validation) IsEmail(field, value string) {
	if !govalidator.IsEmail(value) {
		v.AddError(field, "invalid email address")
	}
}

// IsInt validates if the given value is a valid integer
func (v *Validation) IsInt(field, value string) {
	_, err := strconv.Atoi(value)
	if err != nil {
		v.AddError(field, "this field must be an integer")
	}
}

// IsFloat validates if the given value is a valid floating point number
func (v *Validation) IsFloat(field, value string) {
	_, err := strconv.ParseFloat(value, 64)
	if err != nil {
		v.AddError(field, "this field must be an floating point number")
	}
}

// IsDateISO validates if the given value is a valid ISO date
func (v *Validation) IsDateISO(field, value string) {
	_, err := time.Parse("2006-01-02", value)
	if err != nil {
		v.AddError(field, "this field must be a date in the form of YYYY-MM-DD")
	}
}

// NoSpaces validates if the given value does not contain spaces
func (v *Validation) NoSpaces(field, value string) {
	if govalidator.HasWhitespace(value) {
		v.AddError(field, "whitespaces are not permitted")
	}
}
