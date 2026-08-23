package app

import (
	"errors"
	"fmt"

	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	entranslations "github.com/go-playground/validator/v10/translations/en"
	"github.com/roka-crew/sam2ooh2-api/internal/apperr"
)

type structValidator struct {
	validate *validator.Validate
	trans    ut.Translator
}

func newStructValidator() (*structValidator, error) {
	v := validator.New(validator.WithRequiredStructEnabled())

	// 영어 로케일 및 번역기 설정
	enLocale := en.New()
	uni := ut.New(enLocale, enLocale)
	trans, ok := uni.GetTranslator("en")
	if !ok {
		return nil, fmt.Errorf("failed to get 'en' translator from universal-translator")
	}

	// validator에 영어 번역 규칙 등록
	if err := entranslations.RegisterDefaultTranslations(v, trans); err != nil {
		return nil, fmt.Errorf("failed to register default 'en' translations: %w", err)
	}

	return &structValidator{
		validate: v,
		trans:    trans,
	}, nil
}

func (v *structValidator) Validate(out any) error {
	err := v.validate.Struct(out)

	var validationErrors validator.ValidationErrors
	if errors.As(err, &validationErrors) {
		var errMsgs []string
		for _, e := range validationErrors {
			errMsgs = append(errMsgs, e.Translate(v.trans))
		}
		return apperr.ErrInvalidInput.WithDetails(errMsgs)
	}
	return nil
}
