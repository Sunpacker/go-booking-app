package forms

import (
	"fmt"
	"net/url"
	"testing"
)

func TestForm_Has(t *testing.T) {
	formData := url.Values{}
	formData.Add("a", "123")

	form := getTestForm(formData)

	hasMissing := form.Has("some-field")
	if hasMissing {
		t.Error("form has non-existing field")
	}

	hasA := form.Has("a")
	if !hasA {
		t.Error("field A missing in form")
	}
}

func TestForm_Valid(t *testing.T) {
	form := getTestForm(nil)

	isValid := form.Valid()
	if !isValid {
		t.Error("got invalid when should have been valid")
	}
}

func TestForm_Required(t *testing.T) {
	form := getTestForm(nil)

	form.Required("a", "b", "c")
	if form.Valid() {
		t.Error("form shows valid when required fields missing")
	}

	formData := url.Values{}
	formData.Add("a", "a")
	formData.Add("b", "a")
	formData.Add("c", "a")

	form = getTestForm(formData)
	form.Required("a", "b", "c")
	if !form.Valid() {
		t.Error("shows 'doesn't have required fields' when it 'does'")
	}
}

func TestForm_MinLength(t *testing.T) {
	formData := url.Values{}
	formData.Add("first_name", "Ivan")

	form := getTestForm(formData)

	form.MinLength("first_name", 3)
	if !form.Valid() {
		t.Error("expected 'first_name' value with length 3 to be valid")
	}

	form.MinLength("last_name", 3)
	if form.Valid() {
		t.Error("expected 'last_name' value to be invalid, but its valid")
	}

	isError := form.Errors.Get("last_name")
	if isError == "" {
		t.Error("'last_name' shouldn't be in form, but it is")
	}
}

func TestForm_IsEmail(t *testing.T) {
	emailToTest := "ivan@ya.ru"

	formData := url.Values{}
	formData.Add("email", emailToTest)

	form := getTestForm(formData)

	form.IsEmail("email")
	if !form.Valid() {
		t.Error(fmt.Sprintf("email '%s' should be valid, but got invalid", emailToTest))
	}

	form.IsEmail("non-existing")
	if form.Valid() {
		t.Error("non-existing email must be invalid, but its valid")
	}
}

func getTestForm(formData url.Values) *Form {
	if formData != nil {
		return New(formData)
	}
	return New(url.Values{})
}
