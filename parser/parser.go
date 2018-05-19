package parser

import (
	"bytes"
	"os/user"
	"strings"
	"text/template"

	log "github.com/sirupsen/logrus"
)

const (
	HomePlaceholder           = "HOME"
	UserPlaceholder           = "USER"
	DistroPlaceholder         = "DISTRO"
	PackageManagerPlaceholder = "PCK_MANAGER"

	defaultParserPrefix = "$#"
	defaultParserSuffix = "#$"
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

	parser.parsingData[HomePlaceholder] = usr.HomeDir
	parser.parsingData[UserPlaceholder] = usr.Username

	log.WithField("parsing data", parser.parsingData).
		Debug("parsing data processed")

	return parser, nil
}

func (d *defaultParser) Parse(val string) (string, error) {
	preparedValue := strings.Replace(val, defaultParserPrefix, defaultParserPrefix+".", -1)
	t, err := template.New("").Delims(defaultParserPrefix, defaultParserSuffix).Parse(preparedValue)
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
