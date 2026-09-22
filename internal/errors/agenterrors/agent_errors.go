package agenterrors

import (
	"errors"
	"fmt"
	"strings"
)

type AgentErrorClassification int

const (
	// NonRetriable - операцию не следует повторять
	NonRetriable AgentErrorClassification = iota

	// Retriable - операцию можно повторить
	Retriable
)

var ErrConnectionRefused = errors.New("connection refused")

func Classify(err error) AgentErrorClassification {
	if err == nil {
		return NonRetriable
	}

	if errors.Is(err, ErrConnectionRefused) {
		return Retriable
	}

	return NonRetriable
}

func ClassifyAgentError(err error) error {
	if strings.Contains(err.Error(), "connection refused") {
		return ErrConnectionRefused
	}
	return fmt.Errorf("unknown error %x", err)
}
