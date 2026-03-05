package graphics

import "fmt"

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

func (b *Bitmap) getChunksDimensions(cols, rows int) (chunkWidths, chunkHeights []int) {
	fChunkWidth := float64(b.Width) / float64(cols)
	fChunkHeight := float64(b.Height) / float64(rows)

	chunkWidths = make([]int, cols)
	chunkHeights = make([]int, rows)

	totalWidth := 0
	for col := 0; col < cols-1; col++ {
		w := int(fChunkWidth*float64(col+1) - float64(totalWidth))
		totalWidth += w
		chunkWidths[col] = w
	}
	chunkWidths[cols-1] = b.Width - totalWidth

	totalHeight := 0
	for row := 0; row < rows-1; row++ {
		h := int(fChunkHeight*float64(row+1) - float64(totalHeight))
		totalHeight += h
		chunkHeights[row] = h
	}
	chunkHeights[rows-1] = b.Height - totalHeight

	return
}

func (b *Bitmap) SplitIntoChunks(cols, rows int) ([][]Bitmap, error) {
	if cols >= b.Width || rows >= b.Height {
		return nil, fmt.Errorf("unable to split bitmap into chunks: it is too small")
	}

	chunks := make([][]Bitmap, rows)

	chunkWidths, chunkHeights := b.getChunksDimensions(cols, rows)

	yCursor := 0
	for row := range rows {
		chunks[row] = make([]Bitmap, cols)

		chunkHeight := chunkHeights[row]

		xCursor := 0
		for col := range cols {
			chunkWidth := chunkWidths[col]

			bitmap, err := NewBitmap(chunkWidth, chunkHeight)
			if err != nil {
				return nil, err
			}

			bitmap.Buffer = make([]byte, chunkWidth*chunkHeight)

			for i := range chunkHeight {
				srcStart := yCursor*b.Width + xCursor + i*b.Width
				srcEnd := srcStart + chunkWidth

				dstStart := i * chunkWidth
				dstEnd := dstStart + chunkWidth

				copy(bitmap.Buffer[dstStart:dstEnd], b.Buffer[srcStart:srcEnd])
			}

			chunks[row][col] = *bitmap

			xCursor += chunkWidth
		}

		yCursor += chunkHeight
	}

	return chunks, nil
}
