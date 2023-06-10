package validator

type Validator struct {
}

type IValidater interface {
	Validate() error
}

func (v *Validator) Validate(i interface{}) error {
	if iv, ok := i.(IValidater); ok {
		return iv.Validate()
	}
	return nil
}
