package mock

type Formatter struct{}

func New() Formatter {
	return Formatter{}
}

func (Formatter) Format(tmpl string, _ any) (string, error) {
	return tmpl, nil
}
