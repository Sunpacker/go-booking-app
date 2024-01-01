package forms

import (
	"fmt"
	"github.com/asaskevich/govalidator"
	"net/http"
	"net/url"
	"strings"
)

type Form struct {
	url.Values
	Errors errors
}

func New(data url.Values) *Form {
	return &Form{
		data,
		errors(map[string][]string{}),
	}
}

func (form *Form) Has(field string, r *http.Request) bool {
	fieldFromForm := r.Form.Get(field)
	if fieldFromForm == "" {
		return false
	}
	return true
}

func (form *Form) Valid() bool {
	return len(form.Errors) == 0
}

func (form *Form) Required(fields ...string) {
	for _, field := range fields {
		value := form.Get(field)
		if strings.TrimSpace(value) == "" {
			form.Errors.Add(field, "This field cannot be blank")
		}
	}
}

func (form *Form) MinLength(field string, length int, r *http.Request) bool {
	value := r.Form.Get(field)
	if len(value) < length {
		form.Errors.Add(field, fmt.Sprintf("This field must be at least %d characters long", length))
		return false
	}
	return true
}

func (form *Form) IsEmail(field string) {
	if !govalidator.IsEmail(form.Get(field)) {
		form.Errors.Add(field, "Invalid email address")
	}
}