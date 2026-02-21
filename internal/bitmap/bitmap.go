package bitmap

import "fmt"

type Bitmap struct {
	buffer []byte
	width  int
	height int
}

func NewBitmap(width, height int) (*Bitmap, error) {
	if width <= 0 || height <= 0 {
		return nil, fmt.Errorf("bitmap dimensions (%dx%d) are invalid", width, height)
	}

	return &Bitmap{
		make([]byte, width*height),
		width,
		height,
	}, nil
}
