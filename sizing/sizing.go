package sizing

type OutputSize struct {
	Width  OutputWidth
	Height OutputHeight
}

var OutputAuto = OutputSize{WidthAuto{}, HeightAuto{}}
