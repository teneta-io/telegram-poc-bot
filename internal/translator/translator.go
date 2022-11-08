package translator

import (
	"encoding/json"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"go.uber.org/zap"
	"golang.org/x/text/language"
	"sync"
)

var (
	translator *Translator
	once       sync.Once
)

type Translator struct {
	bundle *i18n.Bundle
}

func NewTranslator() *Translator {
	once.Do(func() {
		bundle := i18n.NewBundle(language.English)
		bundle.RegisterUnmarshalFunc("json", json.Unmarshal)
		bundle.MustLoadMessageFile("./assets/translates/en.json")

		translator = &Translator{
			bundle: bundle,
		}
	})

	return translator
}

func (t *Translator) Translate(msg string, language string, args map[string]interface{}) (str string) {
	var err error

	if str, err = i18n.NewLocalizer(translator.bundle, language).
		Localize(&i18n.LocalizeConfig{MessageID: msg, TemplateData: args}); err != nil {
		zap.S().Error(err.Error())

		return msg
	}

	return str
}
