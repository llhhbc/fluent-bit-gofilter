package parsers

type Parser interface {
	GetName() string
	Parse(src map[string]string) error
}

type ParserConfig struct {
	ParseTimer   []*ParseTimer   `yaml:"parseTimer"`
	ParseRegex   []*ParseRegex   `yaml:"parseRegex"`
	ParseJson    []*ParseJson    `yaml:"parseJson"`
	ParseK8sUid  []*ParseK8sUid  `yaml:"parseK8sUid"`
	ParseK8sName []*ParseK8sName `yaml:"parseK8sName"`
}

var _ Parser = &ParseTimer{}
var _ Parser = &ParseRegex{}
var _ Parser = &ParseJson{}
var _ Parser = &ParseK8sUid{}
var _ Parser = &ParseK8sName{}
