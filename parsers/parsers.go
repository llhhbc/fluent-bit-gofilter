package parsers

type Parser interface {
	GetName() string
	Parse(src map[string]string) error
}

type ParserConfig struct {
	ParseField []ParseField `yaml:"parseField"`
}

type ParseField struct {
	Name string
}

func (t ParseField) GetName() string {
	return t.Name
}

func (t ParseField) Parse(src map[string]string) error {
	return nil
}
