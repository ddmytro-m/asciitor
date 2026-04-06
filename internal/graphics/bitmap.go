package graphics

import (
	"fmt"
	"image"

	"golang.org/x/image/draw"
)

type Bitmap struct {
	Width  int
	Height int
	Buffer []byte
}

func NewBitmap(width, height int) (*Bitmap, error) {
	if width <= 0 || height <= 0 {
		return nil, fmt.Errorf("bitmap dimensions (%dx%d) are invalid", width, height)
	}

	return &Bitmap{
		width,
		height,
		make([]byte, width*height),
	}, nil
}

func NewBitmapFromImage(img image.Image) (*Bitmap, error) {
	bounds := img.Bounds()
	w, h := bounds.Dx(), bounds.Dy()

	var grayImg *image.Gray

	if g, ok := img.(*image.Gray); ok {
		grayImg = g
	} else {
		grayImg = image.NewGray(bounds)
		draw.Draw(grayImg, bounds, img, bounds.Min, draw.Src)
	}

	bitmap, err := NewBitmap(w, h)
	if err != nil {
		return nil, err
	}
	bitmap.Buffer = grayImg.Pix

	return bitmap, nil
}

func (b *Bitmap) FillWithImage(img image.Image) error {
	rect := image.Rect(0, 0, b.Width, b.Height)
	monoImg := image.NewGray(rect)

	draw.BiLinear.Scale(monoImg, rect, img, img.Bounds(), draw.Src, nil)

	b.Buffer = monoImg.Pix

	return nil
}

func (b *Bitmap) GetChunk(x1, y1, x2, y2 int) (*Bitmap, error) {
	if x1 < 0 || y1 < 0 || x2 > b.Width || y2 > b.Height {
		return nil, fmt.Errorf("coordinates out of bounds: [%d, %d] to [%d, %d] for bitmap %dx%d", x1, y1, x2, y2, b.Width, b.Height)
	}

	newWidth := x2 - x1
	newHeight := y2 - y1

	if newWidth <= 0 || newHeight <= 0 {
		return nil, fmt.Errorf("invalid chunk dimensions: %dx%d", newWidth, newHeight)
	}

	newBuffer := make([]byte, newWidth*newHeight)

	for y := range newHeight {
		srcStart := (y1+y)*b.Width + x1
		srcEnd := srcStart + newWidth

		dstStart := y * newWidth

		copy(newBuffer[dstStart:dstStart+newWidth], b.Buffer[srcStart:srcEnd])
	}

	return &Bitmap{
		Width:  newWidth,
		Height: newHeight,
		Buffer: newBuffer,
	}, nil
}
