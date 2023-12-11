package main

import (
	"errors"
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

type BarudakAvengers struct {
	Name  string `validate:"required,maxLen:30"`
	Age   int    `validate:"required,min:25,max:60"`
	Email string `validate:"required,email"`
}

type Validator struct {
}

func (v *Validator) ValidateField(field reflect.Value, validation string) error {
	switch {
	case strings.HasPrefix(validation, "required"):
		if field.IsZero() {
			return errors.New("Tolong diisi")
		}
	case strings.HasPrefix(validation, "maxLen"):
		maxLen, err := strconv.Atoi(strings.TrimPrefix(validation, "maxLen:"))
		if err != nil {
			return err
		}
		if field.Kind() == reflect.String && field.Len() > maxLen {
			return errors.New("Namanya kepanjangan mas")
		}
	case strings.HasPrefix(validation, "minLen"):
		minLen, err := strconv.Atoi(strings.TrimPrefix(validation, "minLen:"))
		if err != nil {
			return err
		}
		if field.Kind() == reflect.String && field.Len() < minLen {
			return errors.New("Gak bisa panjang dikit namanya?")
		}
	case strings.HasPrefix(validation, "min:"):
		min, err := strconv.Atoi(strings.TrimPrefix(validation, "min:"))
		if err != nil {
			return err
		}

		if field.Kind() == reflect.Int && int(field.Int()) < min {
			return errors.New("Kata mama umurnya belom boleh jadi member avengers")
		}

	case strings.HasPrefix(validation, "max:"):
		max, err := strconv.Atoi(strings.TrimPrefix(validation, "max:"))
		if err != nil {
			return err
		}

		if field.Kind() == reflect.Int && int(field.Int()) > max {
			return errors.New("Udah tua, kata mama gak bole masuk avengers")
		}
	case strings.HasPrefix(validation, "email"):
		emailRegex := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
		match, _ := regexp.MatchString(emailRegex, field.String())
		if !match {
			return errors.New("emailnya salah")
		}
	}

	return nil
}

func (v *Validator) Validate(avenger interface{}) error {
	value := reflect.ValueOf(avenger)
	if value.Kind() == reflect.Ptr {
		value = value.Elem()
	}

	for i := 0; i < value.NumField(); i++ {
		field := value.Field(i)
		fieldTag := value.Type().Field(i).Tag.Get("validate")

		validations := strings.Split(fieldTag, ",")

		for _, validation := range validations {
			err := v.ValidateField(field, validation)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func main() {
	validator := Validator{}
	member := BarudakAvengers{
		Name:  "Iron Man",
		Age:   26,
		Email: "tony.stark@gmail.com",
	}

	err := validator.Validate(member)
	if err != nil {
		panic(err)
	}
}
