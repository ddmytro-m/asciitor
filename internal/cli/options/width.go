package options

import (
	"errors"
	"strconv"
	"strings"
)

func validateWidth(s string) error {
	s = strings.Trim(s, " ")
	if len(s) == 0 {
		return errors.New("width not specified")
	}

	if s == "tw" || s == "original" {
		return nil
	}

	matched := rePx.MatchString(s)
	if matched {
		pxStr := s[:len(s)-2] // 1234|px
		if pxStr[0] == '0' {
			return errors.New("leading zeros are forbidden")
		}
		pxNum, err := strconv.Atoi(pxStr)
		if err != nil {
			return err
		} else if pxNum <= 0 {
			return errors.New("width must be a positive integer")
		}
		return nil
	}

	matched = reCols.MatchString(s)
	if matched {
		countStr := s[:len(s)-1] // 1234|M
		if countStr[0] == '0' {
			return errors.New("leading zeros are forbidden")
		}
		countNum, err := strconv.Atoi(countStr)
		if err != nil {
			return err
		} else if countNum <= 0 {
			return errors.New("characters count must be a positive integer")
		}
		return nil
	}

	return errors.New("unknown width value")
}
