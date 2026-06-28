package options

import (
	"errors"

	"github.com/ddmytro-m/asciitor/internal/graphics"
)

const maxBlockSize = graphics.DensityMapMaxBlockSize

func validateBlockSize(val int) error {
	if val <= 0 {
		return errors.New("block size must be positive")
	}

	if val > maxBlockSize {
		return errors.New("block size is too big")
	}

	return nil
}

// idk, seems reasonable
const maxFontSize = (1 << 15) - 1

func validateFontSize(val int) error {
	if val <= 0 {
		return errors.New("font size must be positive")
	}

	if val > maxFontSize {
		return errors.New("font size is too big")
	}

	return nil
}

func validateFaceIndex(val int) error {
	if val < 0 {
		return errors.New("font face must be positive")
	}

	return nil
}
