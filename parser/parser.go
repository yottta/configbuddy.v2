package parser

import (
	"bytes"
	"os/user"
	"strings"
	"text/template"

	log "github.com/sirupsen/logrus"
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
	parser := &defaultParser{
		parsingData: make(map[string]string),
	}

	usr, err := user.Current()
	if err != nil {
		return nil, err
	}

	parser.parsingData[homePlaceholder] = usr.HomeDir
	parser.parsingData[userPlaceholder] = usr.Username

	log.WithField("parsing data", parser.parsingData).
		Debug("parsing data processed")

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
