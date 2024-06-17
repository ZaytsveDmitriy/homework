package hw09structvalidator

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
)

var ErrNotStruct = errors.New("arg not a struct")

type ValidationError struct {
	Field string
	Err   error
}

type ValidationErrors []ValidationError

func (v ValidationErrors) Error() string {
	res := strings.Builder{}
	for _, val := range v {
		res.WriteString(fmt.Sprintf("field: %s, error: %s\n", val.Field, val.Err))
	}
	return res.String()
}

func Validate(v interface{}) error {
	// Place your code here.
	var resError ValidationErrors
	refl := reflect.ValueOf(v)
	if refl.Type().Kind() != reflect.Struct {
		return ErrNotStruct
	}

	for i := 0; i < refl.NumField(); i++ {
		var err error

		tag, has := refl.Type().Field(i).Tag.Lookup("validate")
		if !has {
			continue
		}

		if !refl.Type().Field(i).IsExported() {
			continue
		}

		fName := refl.Type().Field(i).Name
		switch refl.Field(i).Type().Kind() { //nolint:exhaustive
		case reflect.Int:
			val := refl.Field(i).Int()
			err = ValidateInt(tag, val)
		case reflect.String:
			val := refl.Field(i).String()
			err = ValidateString(tag, val)
		case reflect.Slice:
			l := refl.Field(i).Len()
			if l < 1 {
				continue
			}

			switch refl.Field(i).Type().Elem().Kind() { //nolint:exhaustive
			case reflect.Int:
				val := refl.Field(i).Slice(0, l)
				data := make([]int64, 0, l)
				for i := 0; i < l; i++ {
					data = append(data, val.Index(i).Int())
				}
				err = ValidateInt(tag, data...)

			case reflect.String:
				val := refl.Field(i).Slice(0, l)

				data := make([]string, 0, l)
				for i := 0; i < l; i++ {
					data = append(data, val.Index(i).String())
				}
				err = ValidateString(tag, data...)
			default:
				continue
			}
		default:
			continue
		}

		if err != nil {
			if errors.Is(err, ErrTemplateInvalid) {
				return err
			}

			valErr := ValidationError{fName, err}
			resError = append(resError, valErr)
		}
	}

	if len(resError) > 0 {
		return resError
	}
	return nil
}

func ValidateInt(pattern string, args ...int64) error {
	var err error
	filter, err := NewIntFilter(pattern)
	if err != nil {
		return err
	}

	err = filter.Validate(args...)
	if err != nil {
		return err
	}

	return nil
}

func ValidateString(pattern string, args ...string) error {
	var err error
	filter, err := NewStringFilter(pattern)
	if err != nil {
		return err
	}

	err = filter.Validate(args...)
	if err != nil {
		return err
	}

	return nil
}
