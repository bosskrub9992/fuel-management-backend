package validators

import (
	"errors"
	"strings"

	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	entranslations "github.com/go-playground/validator/v10/translations/en"
)

type RequestValidator struct {
	validator *validator.Validate
	trans     ut.Translator
}

func NewRequestValidator() *RequestValidator {
	v := validator.New()
	en := en.New()
	uni := ut.New(en, en)
	translator, found := uni.GetTranslator("en")
	if !found {
		panic("translator not found")
	}

	if err := entranslations.RegisterDefaultTranslations(v, translator); err != nil {
		panic(err)
	}

	return &RequestValidator{
		validator: v,
		trans:     translator,
	}
}

func (rv *RequestValidator) Validate(s any) error {
	if err := rv.validator.Struct(s); err != nil {
		fieldErrors := err.(validator.ValidationErrors)
		messages := []string{}
		for _, e := range fieldErrors {
			messages = append(messages, e.Translate(rv.trans))
		}
		return errors.New(strings.Join(messages, ", "))
	}

	return nil
}
