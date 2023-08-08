package forms

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/asaskevich/govalidator"
)

// Form creates a custom form struct
type Form struct {
	url.Values
	Errors errors
}

func (f *Form) Valid() bool {
	return len(f.Errors) == 0
}

// New initializes the form struct
func New(data url.Values) *Form {
	return &Form{
		data,
		errors(map[string][]string{}),
	}
}

func (f *Form) Has(field string, r *http.Request) bool {
	x := r.Form.Get(field)
	return x != ""
	/* if x == "" {
		f.Errors.Add(field, "This field cannot be empty")
		return false
	}

	return true */
}

// Reqired checks for required fields
func (f *Form) Required(fields ...string) {
	for _, field := range fields {
		value := f.Get(field)

		if strings.TrimSpace(value) == "" {
			f.Errors.Add(field, "This field is required")
		}
	}
}

// MinLength checks string for minimum length
func (f *Form) MinLength(field string, length int, r *http.Request) bool {
	x := r.Form.Get(field)

	if len(x) < length {
		f.Errors.Add(field, fmt.Sprintf("This field must be at least %d characters long", length))
		return false
	}

	return true
}

// IsEmail checks for valid email address
func (f *Form) IsEmail(field string) {
	if !govalidator.IsEmail(f.Get(field)) {
		f.Errors.Add(field, "Invalid email address")
	}
}

// PasswordsMatch checks that passwords are equal
func (f *Form) PasswordsMatch(field1 string, field2 string, r *http.Request) {
	p1 := r.Form.Get(field1)
	p2 := r.Form.Get(field2)

	if strings.Compare(p1, p2) != 0 || p1 != p2 || len(p1) != len(p2) {
		f.Errors.Add(field1, "Passwords don't match")
		f.Errors.Add(field2, "Passwords don't match")
	}

}
