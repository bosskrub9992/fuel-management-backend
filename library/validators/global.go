package validators

var Validator *RequestValidator

func init() {
	Validator = NewRequestValidator()
}

func Validate(s any) error {
	return Validator.Validate(s)
}
