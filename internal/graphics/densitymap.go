package graphics

import (
	"fmt"
	"math"
)

type DensityMap struct {
	width, height int
	cells         []uint64
}

const maxBlockSize = (64 - 8) / 2

func NewDensityMapFromBitmap(bitmap *Bitmap, blockSize int) (*DensityMap, error) {
	if bitmap.Buffer == nil {
		return nil, fmt.Errorf("cannot convert an empty buffer into a density map")
	}
	if blockSize > maxBlockSize {
		return nil, fmt.Errorf("given chunk size is too big (max value is 2^%d).", maxBlockSize)
	}

	width := (bitmap.Width + blockSize - 1) / blockSize
	height := (bitmap.Height + blockSize - 1) / blockSize
	cells := make([]uint64, width*height)

	for y := range bitmap.Height {
		for x := range bitmap.Width {
			targetCell := (y/blockSize)*width + (x / blockSize)
			cells[targetCell] += uint64(bitmap.Buffer[y*bitmap.Width+x])
		}
	}

	return &DensityMap{width, height, cells}, nil
}

func (dm *DensityMap) GetDistance(another *DensityMap) (float64, error) {
	if dm.width <= 0 || dm.height <= 0 {
		return -1, fmt.Errorf("cannot get distance of empty density maps")
	}

	if dm.width != another.width || dm.height != another.height {
		return -1, fmt.Errorf("density maps are not comparable: (%dx%d) vs (%dx%d)", dm.width, dm.height, another.width, another.height)
	}

	var sumOfSquares float64

	for i := range dm.cells {
		// convert to float64 to avoid overflow
		diff := float64(dm.cells[i]) - float64(another.cells[i])
		sumOfSquares += diff * diff
	}

	return math.Sqrt(sumOfSquares), nil
}
