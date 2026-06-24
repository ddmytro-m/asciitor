package options

import (
	"errors"
	"strings"
)

func validateOutput(s string) error {
	s = strings.Trim(s, " ")
	if len(s) == 0 {
		return errors.New("output not specified")
	}

	return nil
}
