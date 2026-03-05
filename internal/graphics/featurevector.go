package graphics

import (
	"fmt"
	"math"
)

type FeatureVector struct {
	width, height int
	values        []byte
}

func NewFeatureVectorFromBitmap(bitmap Bitmap, width, height int) (*FeatureVector, error) {
	if width <= 0 || height <= 0 {
		return nil, fmt.Errorf("feature vector must have positive dimensions")
	}

	chunks, err := bitmap.SplitIntoChunks(width, height)
	if err != nil {
		return nil, err
	}

	fv := FeatureVector{width, height, make([]byte, width*height)}

	for i := range len(chunks) {
		for j := range len(chunks[0]) {
			chunk := chunks[i][j]

			if len(chunk.Buffer) == 0 {
				return nil, fmt.Errorf("cannot create feature vector value from an empty buffer")
			}

			var sum int
			for _, p := range chunk.Buffer {
				sum += int(p)
			}
			fv.values[i*width+j] = byte(sum / len(chunk.Buffer))
		}
	}

	return &fv, nil
}

func (f *FeatureVector) GetDistance(another *FeatureVector) (float64, error) {
	if f.width <= 0 || f.height <= 0 {
		return -1, fmt.Errorf("cannot get distance of empty feature vectors")
	} else if f.width != another.width || f.height != another.height {
		return -1, fmt.Errorf("feature vector of dimensions (%dx%d) is not comparable with vector of dimensions (%dx%d)", f.width, f.height, another.width, another.height)
	}

	var d2 float64
	for i := range len(f.values) {
		diff := float64(f.values[i]) - float64(another.values[i])
		d2 += diff * diff
	}

	return math.Sqrt(d2), nil
}
