// Package validate checks that a request says what it has to say, before a
// handler acts on it.
//
// The rules are struct tags, read by go-playground/validator:
//
//	type loginRequest struct {
//		Username string `json:"username" validate:"required"`
//	}
//
// What this package adds is the answer. A failed rule comes back as a
// respond.Fault carrying a 400 and a sentence naming the field the way the
// request named it — "email must be an email address", not "Key: 'Email'
// Error:Field validation for 'Email' failed on the 'email' tag".
package validate

import (
	"errors"
	"fmt"
	"net/http"
	"reflect"
	"strings"

	"github.com/go-playground/validator/v10"

	"xermess/internal/api/http/respond"
)

// instance is shared: building one is expensive, and it is safe to use from
// every request at once.
var instance = newValidator()

// messages is what to say about each rule, added to the field's name. The
// built-in rules are listed here; a package can add its own with Register.
var messages = map[string]string{
	"required":   "is required",
	"email":      "must be an email address",
	"max":        "is too long",
	"min":        "is too short",
	"oneof":      "is not one of the values allowed",
	"startswith": "does not start the way it has to",
}

func newValidator() *validator.Validate {
	v := validator.New(validator.WithRequiredStructEnabled())

	// Name a field in an error the way the request named it, so the panel can
	// show the message next to the box someone typed in.
	v.RegisterTagNameFunc(func(field reflect.StructField) string {
		name := strings.Split(field.Tag.Get("json"), ",")[0]
		if name == "" || name == "-" {
			return field.Name
		}

		return name
	})

	return v
}

// Register adds a rule of our own, for something the library cannot know
// about — what a field name may look like, which types exist. The message is
// what to say when the rule is broken, after the field's name.
//
// Call it from the init of the package whose rule it is, so the rule and the
// thing it describes stay together.
func Register(tag, message string, valid func(value string) bool) {
	if err := instance.RegisterValidation(tag, func(fl validator.FieldLevel) bool {
		return valid(fl.Field().String())
	}); err != nil {
		// Only a bad tag name gets here, which is a mistake in our own code
		// rather than anything a request can cause.
		panic(fmt.Sprintf("validate: register %q: %v", tag, err))
	}

	messages[tag] = message
}

// Struct checks a request against its tags. It returns a respond.Fault, so a
// handler passes whatever comes back to respond.Failure and is done.
func Struct(v any) error {
	err := instance.Struct(v)
	if err == nil {
		return nil
	}

	var broken validator.ValidationErrors
	if !errors.As(err, &broken) {
		// The value was not a struct at all: our mistake, not the caller's.
		return err
	}

	return respond.Fault{Status: http.StatusBadRequest, Message: Message(broken[0])}
}

// Message is the sentence for one broken rule. Only the first is reported:
// the panel shows one line, and a form is easier to fix one thing at a time.
func Message(broken validator.FieldError) string {
	said, ok := messages[broken.Tag()]
	if !ok {
		said = "is not valid"
	}

	// A length rule says which length it wanted; the others speak for
	// themselves.
	if param := broken.Param(); param != "" {
		switch broken.Tag() {
		case "max":
			said = fmt.Sprintf("must be at most %s characters", param)
		case "min":
			said = fmt.Sprintf("must be at least %s characters", param)
		case "oneof":
			said = "must be one of: " + strings.ReplaceAll(param, " ", ", ")
		case "startswith":
			said = fmt.Sprintf("must start with %q", param)
		}
	}

	return broken.Field() + " " + said
}
