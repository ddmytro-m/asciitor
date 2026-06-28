package testutils

import (
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"os"
)

func PngToImage(path string) (image.Image, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("could not load image: %v", err)
	}

	img, err := png.Decode(f)
	if err != nil {
		return nil, fmt.Errorf("could not decode image: %v", err)
	}

	return img, nil
}

func JpegToImage(path string) (image.Image, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("could not load image: %v", err)
	}

	img, err := jpeg.Decode(f)
	if err != nil {
		return nil, fmt.Errorf("could not decode image: %v", err)
	}

	return img, nil
}
