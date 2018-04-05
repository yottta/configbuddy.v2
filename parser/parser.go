package parser

import (
	"bytes"
	"os/user"
	"strings"
	"text/template"
)

const (
	homePlaceholder           = "HOME"
	userPlaceholder           = "USER"
	distroPlaceholder         = "DISTRO"
	packageManagerPlaceholder = "PCK_MANAGER"
)

type Parser interface {
	Parse(val string) (string, error)
}

type defaultParser struct {
	parsingData map[string]string
}

func NewParser() (Parser, error) {
	parser := &defaultParser{}

	dat := make(map[string]string)

	usr, err := user.Current()
	if err != nil {
		return nil, err
	}

	dat[homePlaceholder] = usr.HomeDir
	dat[userPlaceholder] = usr.Username

	parser.parsingData = dat

	return parser, nil
}

func (d *defaultParser) Parse(val string) (string, error) {
	t, err := template.New("").Delims("$#", "#$").Parse(strings.Replace(val, "$#", "$#.", -1))
	if err != nil {
		return "", err
	}

	var bytes bytes.Buffer
	err = t.Execute(&bytes, d.parsingData)
	if err != nil {
		return "", err
	}
	return bytes.String(), nil
}
