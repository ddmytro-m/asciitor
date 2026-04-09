package graphics

type SummedAreaTable struct {
	width, height int
	data          []uint64
}

func NewSummedAreaTable(bitmap *Bitmap) *SummedAreaTable {
	width, height := bitmap.Width, bitmap.Height
	stride := width + 1
	data := make([]uint64, stride*(height+1))

	for y := range height {
		var rowSum uint64 = 0
		rowStart := y * width

		for x := range width {
			pixel := uint64(bitmap.Buffer[rowStart+x])
			rowSum += pixel
			data[(y+1)*stride+(x+1)] = data[y*stride+(x+1)] + rowSum
		}
	}
	return &SummedAreaTable{width, height, data}
}

func (sat *SummedAreaTable) GetSum(x1, y1, x2, y2 int) uint64 {
	if x1 < 0 {
		x1 = 0
	}
	if y1 < 0 {
		y1 = 0
	}
	if x2 > sat.width {
		x2 = sat.width
	}
	if y2 > sat.height {
		y2 = sat.height
	}

	if x2 <= x1 || y2 <= y1 {
		return 0
	}

	stride := sat.width + 1

	A := sat.data[y1*stride+x1]
	B := sat.data[y1*stride+x2]
	C := sat.data[y2*stride+x1]
	D := sat.data[y2*stride+x2]

	return D - B - C + A
}

func (sat *SummedAreaTable) ExtractDensityMap(x, y, glyphW, glyphH, chunkSize int) *DensityMap {
	dmWidth := (glyphW + chunkSize - 1) / chunkSize
	dmHeight := (glyphH + chunkSize - 1) / chunkSize
	cells := make([]uint64, dmWidth*dmHeight)

	for cy := range dmHeight {
		for cx := range dmWidth {
			x1 := x + (cx * chunkSize)
			y1 := y + (cy * chunkSize)

			x2 := x1 + chunkSize
			y2 := y1 + chunkSize

			if x2 > x+glyphW {
				x2 = x + glyphW
			}
			if y2 > y+glyphH {
				y2 = y + glyphH
			}

			cells[cy*dmWidth+cx] = sat.GetSum(x1, y1, x2, y2)
		}
	}

	return &DensityMap{
		width:  dmWidth,
		height: dmHeight,
		cells:  cells,
	}
}
