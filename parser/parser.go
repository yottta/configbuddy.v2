package parser

import (
	"bytes"
	log "log/slog"
	"os/user"
	"strings"
	"text/template"

	"github.com/Knetic/govaluate"
)

const (
	HomePlaceholder           = defaultParserSuffix + "HOME" + defaultParserSuffix
	UserPlaceholder           = defaultParserSuffix + "USER" + defaultParserSuffix
	DistroPlaceholder         = defaultParserSuffix + "DISTRO" + defaultParserSuffix
	PackageManagerPlaceholder = defaultParserSuffix + "PCK_MANAGER" + defaultParserSuffix

	HomePlaceholderName           = "HOME"
	UserPlaceholderName           = "USER"
	DistroPlaceholderName         = "DISTRO"
	PackageManagerPlaceholderName = "PCK_MANAGER"

	defaultParserPrefix = "$#"
	defaultParserSuffix = "#$"
)

type Parser interface {
	Parse(val string) (string, error)
	EvaluateCondition(condition string) (bool, error)
}

type defaultParser struct {
	parsingData    map[string]string
	conditionsData map[string]interface{}
}

func NewParser() (Parser, error) {
	parser := &defaultParser{
		parsingData:    make(map[string]string),
		conditionsData: make(map[string]interface{}),
	}

	usr, err := user.Current()
	if err != nil {
		return nil, err
	}

	parser.parsingData[HomePlaceholderName] = usr.HomeDir
	parser.parsingData[UserPlaceholderName] = usr.Username
	packageManagerName, err := PckManager()
	if err != nil {
		return nil, err
	}
	parser.parsingData[PackageManagerPlaceholderName] = packageManagerName

	for k, v := range parser.parsingData {
		newKey := strings.TrimRight(strings.TrimLeft(k, defaultParserPrefix), defaultParserSuffix)
		parser.conditionsData[newKey] = v
	}
	log.With("parsing data", parser.parsingData).
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

func (d *defaultParser) EvaluateCondition(condition string) (bool, error) {
	expression, err := govaluate.NewEvaluableExpression(condition)
	if err != nil {
		return false, err
	}

	result, err := expression.Evaluate(d.conditionsData)
	if err != nil {
		return false, err
	}

	res, ok := result.(bool)
	if !ok {
		log.With("condition", condition).Warn("failed to evaluate the condition")
		return false, nil
	}
	return res, nil
}
