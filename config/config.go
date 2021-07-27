package config

import (
	"errors"
	"fmt"
	"strings"
	"unicode"
)

const (
	stateStartProperty = iota
	stateStartName
	stateExpectEqual
	stateStartValue
)

const flagProperty = "-P"

func tokenizeCommandLine(commandLine string) []string {
	tokens := make([]string,0)
	currentToken := ""
	for _, char := range commandLine {
		if unicode.IsSpace(char) {
			if len(currentToken) > 0 {
				tokens = append(tokens, currentToken)
				currentToken = ""
			}
			continue
		}
		if string(char) == "=" {
			if len(currentToken) > 0 {
				tokens = append(tokens, currentToken)
			}
			tokens = append(tokens, "=")
			currentToken = ""
			continue
		}
		currentToken = currentToken + string(char)
	}
	if len(currentToken) > 0 {
		tokens = append(tokens, currentToken)
	}
	return tokens
}

func parseTokens(tokens []string) (map[string]string, error) {
	properties := make(map[string]string)
	property := ""
	state := stateStartProperty
	for index, token := range tokens {
		switch state {
		case stateStartProperty:
			if token == flagProperty {
				state = stateStartName
			} else {
				return nil, errors.New(fmt.Sprintf("expected '-P' at position %d", index))
			}
		case stateStartName:
			property = token
			state = stateExpectEqual
		case stateExpectEqual:
			if token != "=" {
				return nil, errors.New(fmt.Sprintf("expected '=' at position %d", index))
			}
			state = stateStartValue
		case stateStartValue:
			properties[property] = token
			state = stateStartProperty
		}
	}
	return properties, nil
}

func ParseCommandLine(args []string, params map[string]string) error {
	if len(args) <= 1 {
		return nil
	}
	commandLine := strings.Join(args[1:], " ")
	tokens := tokenizeCommandLine(commandLine)
	parsed, err := parseTokens(tokens)
	if err != nil {
		return err
	}
	for key, value := range parsed {
		params[key] = value
	}
	return nil
}
