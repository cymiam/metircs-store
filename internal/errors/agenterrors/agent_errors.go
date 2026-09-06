package agenterrors

import (
	"errors"
	"strings"
)

type AgentErrorClassification int

const (
	// NonRetriable - операцию не следует повторять
	NonRetriable AgentErrorClassification = iota

	// Retriable - операцию можно повторить
	Retriable
)

var ConnectionRefusedError = errors.New("connection refused")

func Classify(err error) AgentErrorClassification {
	if err == nil {
		return NonRetriable
	}

	if errors.Is(err, ConnectionRefusedError) {
		return Retriable
	}

	return NonRetriable
}

func ClassifyAgentError(err error) error {
	if strings.Contains(err.Error(), "connection refused") {
		return ConnectionRefusedError
	}
	return errors.New("unknown error")
}
