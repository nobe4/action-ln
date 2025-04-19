package format

type Formatter interface {
	Format(tmpl string, data any) (string, error)
}
