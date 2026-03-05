package testutils

import (
	"fmt"
	"image"
	"image/draw"
	"image/jpeg"
	"image/png"
	"os"

	"github.com/ddmytro-m/asciitor/internal/graphics"
)

func ImgToBitmap(img image.Image) (*graphics.Bitmap, error) {
	bounds := img.Bounds()
	w, h := bounds.Dx(), bounds.Dy()

	monoImg := image.NewGray(img.Bounds())
	draw.Draw(monoImg, bounds, img, bounds.Min, draw.Src)
	buffer := monoImg.Pix

	bitmap, err := graphics.NewBitmap(w, h)
	if err != nil {
		return nil, err
	}
	bitmap.Buffer = buffer

	return bitmap, nil
}

func PngToBitmap(path string) (*graphics.Bitmap, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("could not load image: %v", err)
	}

	img, err := png.Decode(f)
	if err != nil {
		return nil, fmt.Errorf("could not decode image: %v", err)
	}

	return ImgToBitmap(img)
}

func JpegToBitmap(path string) (*graphics.Bitmap, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("could not load image: %v", err)
	}

	img, err := jpeg.Decode(f)
	if err != nil {
		return nil, fmt.Errorf("could not decode image: %v", err)
	}

	return ImgToBitmap(img)
}
