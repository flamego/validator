package validator

import (
	"bytes"
	"fmt"
	"reflect"
	"strings"
)

const fieldErrMsg = "Key: %q Error: Field validation for %q failed on the %q tag"

// ValidationErrorsTranslations is the translation return type
type ValidationErrorsTranslations map[string]string

// InvalidValidationError describes an invalid argument passed to
// `Struct`, `StructExcept`, StructPartial` or `Field`
type InvalidValidationError struct {
	Type reflect.Type
}

// Error returns InvalidValidationError message
func (e *InvalidValidationError) Error() string {
	if e.Type == nil {
		return "validator: (nil)"
	}
	return "validator: (nil " + e.Type.String() + ")"
}

// ValidationErrors is an array of FieldError's
// for use in custom error messages post validation.
type ValidationErrors []FieldError

// Error is intended for use in development + debugging and not intended to be a production error message.
// It allows ValidationErrors to subscribe to the Error interface.
// All information to create an error message specific to your application is contained within
// the FieldError found within the ValidationErrors array
func (ve ValidationErrors) Error() string {
	var buf bytes.Buffer
	var fe *fieldError
	for i := 0; i < len(ve); i++ {
		fe = ve[i].(*fieldError)
		buf.WriteString(fe.Error())
		buf.WriteString("\n")
	}
	return strings.TrimSpace(buf.String())
}

// FieldError contains all functions to get error details
type FieldError interface {
	// Tag returns the validation tag that failed. if the
	// validation was an alias, this will return the
	// alias name and not the underlying tag that failed.
	//
	// eg. alias "iscolor": "hexcolor|rgb|rgba|hsl|hsla"
	// will return "iscolor"
	Tag() string

	// ActualTag returns the validation tag that failed, even if an
	// alias the actual tag within the alias will be returned.
	// If an 'or' validation fails the entire or will be returned.
	//
	// eg. alias "iscolor": "hexcolor|rgb|rgba|hsl|hsla"
	// will return "hexcolor|rgb|rgba|hsl|hsla"
	ActualTag() string

	// Namespace returns the namespace for the field error, with the tag
	// name taking precedence over the field's actual name.
	//
	// eg. JSON name "User.fname"
	//
	// See StructNamespace() for a version that returns actual names.
	//
	// NOTE: this field can be blank when validating a single primitive field
	// using validate.Field(...) as there is no way to extract it's name
	Namespace() string

	// StructNamespace returns the namespace for the field error, with the field's
	// actual name.
	//
	// eq. "User.FirstName" see Namespace for comparison
	//
	// NOTE: this field can be blank when validating a single primitive field
	// using validate.Field(...) as there is no way to extract its name
	StructNamespace() string

	// Field returns the fields name with the tag name taking precedence over the
	// field's actual name.
	//
	// eq. JSON name "fname"
	// see StructField for comparison
	Field() string

	// StructField returns the field's actual name from the struct, when able to determine.
	//
	// eq.  "FirstName"
	// see Field for comparison
	StructField() string

	// Value returns the actual field's value in case needed for creating the error
	// message
	Value() interface{}

	// Param returns the param value, in string form for comparison; this will also
	// help with generating an error message
	Param() string

	// Kind returns the Field's reflect Kind
	//
	// eg. time.Time's kind is a struct
	Kind() reflect.Kind

	// Type returns the Field's reflect Type
	//
	// eg. time.Time's type is time.Time
	Type() reflect.Type

	// Error returns the FieldError's message
	Error() string
}

// compile time interface checks
var _ FieldError = new(fieldError)
var _ error = new(fieldError)

// fieldError contains a single field's validation error along
// with other properties that may be needed for error message creation
// it complies with the FieldError interface
type fieldError struct {
	v              *Validate
	tag            string
	actualTag      string
	ns             string
	structNs       string
	fieldLen       uint8
	structfieldLen uint8
	value          interface{}
	param          string
	kind           reflect.Kind
	typ            reflect.Type
}

// Tag returns the validation tag that failed.
func (fe *fieldError) Tag() string {
	return fe.tag
}

// ActualTag returns the validation tag that failed, even if an
// alias the actual tag within the alias will be returned.
func (fe *fieldError) ActualTag() string {
	return fe.actualTag
}

// Namespace returns the namespace for the field error, with the tag
// name taking precedence over the field's actual name.
func (fe *fieldError) Namespace() string {
	return fe.ns
}

// StructNamespace returns the namespace for the field error, with the field's
// actual name.
func (fe *fieldError) StructNamespace() string {
	return fe.structNs
}

// Field returns the field's name with the tag name taking precedence over the
// field's actual name.
func (fe *fieldError) Field() string {
	return fe.ns[len(fe.ns)-int(fe.fieldLen):]
}

// StructField returns the field's actual name from the struct, when able to determine.
func (fe *fieldError) StructField() string {
	// return fe.structField
	return fe.structNs[len(fe.structNs)-int(fe.structfieldLen):]
}

// Value returns the actual field's value in case needed for creating the error
// message
func (fe *fieldError) Value() interface{} {
	return fe.value
}

// Param returns the param value, in string form for comparison; this will
// also help with generating an error message
func (fe *fieldError) Param() string {
	return fe.param
}

// Kind returns the Field's reflect Kind
func (fe *fieldError) Kind() reflect.Kind {
	return fe.kind
}

// Type returns the Field's reflect Type
func (fe *fieldError) Type() reflect.Type {
	return fe.typ
}

// Error returns the fieldError's error message
func (fe *fieldError) Error() string {
	return fmt.Sprintf(fieldErrMsg, fe.ns, fe.Field(), fe.tag)
}
