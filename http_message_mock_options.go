package httpmessagemock

type OptionFunc func(*MessageMock)

func WithTranslatorOption(translator TranslatorAware) OptionFunc {
	return func(m *MessageMock) {
		m.translator = translator
	}
}
