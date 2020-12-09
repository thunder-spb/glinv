package forms

import (
	"fmt"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"unicode/utf8"
)

// EmailRX regular expression for sanity checking the format of an email address.
// Patternr recommended by the W3C and Web Hypertext Application Technology Working Group on 21.02.2020
var EmailRX = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+\\/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")

// Form struct, which anonymously embeds a url.Values object
// (to hold the form data) and an Errors field to hold any validation errors
// for the form data.
type Form struct {
	url.Values
	Errors errors
}

// New function to initialize a custom Form struct. Notice that
// this takes the form data as the parameter?
func New(data url.Values) *Form {
	return &Form{
		data,
		errors(map[string][]string{}),
	}
}

// Required method to check that specific fields in the form
// data are present and not blank. If any fields fail this check, add the
// appropriate message to the form errors.
func (f *Form) Required(fields ...string) {
	for _, field := range fields {
		value := f.Get(field)
		if strings.TrimSpace(value) == "" {
			f.Errors.Set(field, "This field cannot be blank")
		}
	}
}

// MaxLength method to check that a specific field in the form
// contains a maximum number of characters. If the check fails then add the
// appropriate message to the form errors.
func (f *Form) MaxLength(field string, l int) {
	value := f.Get(field)
	if value == "" {
		return
	}

	if utf8.RuneCountInString(value) > l {
		f.Errors.Set(field, fmt.Sprintf("This field is too long (maximum is %d characters)", l))
	}
}

// PermittedPorts method to check that a specific field in the form
// matches one of a set of specific permitted values. If the check fails
// then add the appropriate message to the form errors.
func (f *Form) PermittedPorts(field string, portMin int, portMax int) {
	value := f.Get(field)
	if value == "" {
		return
	}

	port, _ := strconv.Atoi(value)
	if port <= portMin || port >= portMax {
		f.Errors.Set(field, "Value must be between 0 and 65535")
	}
}

// MinLength method to check that a specific field in the form
// contains a minimum number of characters. If the check fails then add the
// appropriate message to the form errors.
func (f *Form) MinLength(field string, d int) {
	value := f.Get(field)
	if value == "" {
		return
	}
	if utf8.RuneCountInString(value) < d {
		f.Errors.Set(field, fmt.Sprintf("This field is too short (minimum is %d characters)", d))
	}
}

// MatchesPattern method to check that a specific field in the form
// matches a regular expression. If the check fails then add the
// appropriate message to the form errors.
func (f *Form) MatchesPattern(field string, pattern *regexp.Regexp) {
	value := f.Get(field)
	if value == "" {
		return
	}
	if !pattern.MatchString(value) {
		f.Errors.Set(field, "This field is invalid")
	}
}

// Valid method which returns true if there are no errors.
func (f *Form) Valid() bool {
	return len(f.Errors) == 0
}
