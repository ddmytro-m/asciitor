package freetype

//#cgo pkg-config: freetype2
//#include "freetype.h"
import "C"
import (
	"fmt"
	"runtime"
	"unsafe"
)

func getErrorMessage(errorCode int) string {
	switch errorCode {
	case C.FT_Err_Invalid_Character_Code:
		return "invalid character code"
	case C.FT_Err_Invalid_Glyph_Index:
		return "no glyph present"
	case C.FT_Err_Cannot_Render_Glyph:
		return "cannot render glyph"
	case C.FT_Err_Invalid_Pixel_Size:
		return "invalid font size"
	case C.FT_Err_Out_Of_Memory:
		return "out of memory"
	default:
		return "unknown error"
	}
}

type FontParams struct {
	Buffer []byte
}

type FontProperties struct {
	FamilyName  string
	FacesAmount int
}

func getCFontParams(params FontParams, p *runtime.Pinner) C.FontParams {
	p.Pin(&params.Buffer[0])

	return C.FontParams{
		buffer:     (*C.uchar)(unsafe.Pointer(&params.Buffer[0])),
		bufferSize: C.int(len(params.Buffer)),
	}
}

func GetFontProperties(params FontParams) (FontProperties, error) {
	if len(params.Buffer) == 0 {
		return FontProperties{}, fmt.Errorf("invalid font buffer: buffer is empty")
	}

	var p runtime.Pinner
	defer p.Unpin()

	cFontParams := getCFontParams(params, &p)
	var cProperties C.FontProperties

	errCode := C.getFont(cFontParams, &cProperties)
	if errCode != 0 {
		return FontProperties{}, fmt.Errorf("%s (error code: 0x%02x)", getErrorMessage(int(errCode)), errCode)
	}
	defer C.freeFont(cProperties)

	return FontProperties{
		FamilyName:  C.GoString(cProperties.familyName),
		FacesAmount: int(cProperties.facesAmount),
	}, nil
}

func GetFaces(params FontParams) ([]FaceProperties, error) {
	if len(params.Buffer) == 0 {
		return []FaceProperties{}, fmt.Errorf("invalid font buffer: buffer is empty")
	}

	var p runtime.Pinner
	defer p.Unpin()

	cFontParams := getCFontParams(params, &p)

	var cFaces *C.FaceProperties
	var cLength C.int

	errCode := C.getFontFaces(cFontParams, &cFaces, &cLength)
	if errCode != 0 {
		return nil, fmt.Errorf("%s (error code: 0x%02x)", getErrorMessage(int(errCode)), errCode)
	}

	length := int(cLength)
	defer C.freeFaces(cFaces, cLength)

	cSlice := unsafe.Slice(cFaces, length)

	goFaces := make([]FaceProperties, length)
	for i := range length {
		goFaces[i] = FaceProperties{
			Index:     int(cSlice[i].index),
			StyleName: C.GoString(cSlice[i].styleName),
			Monospace: bool(cSlice[i].monospace),
		}
	}

	return goFaces, nil
}

type FaceParams struct {
	FontParams FontParams
	FaceIndex  int
}

type FaceProperties struct {
	Index     int
	StyleName string
	Monospace bool
}

func getCFaceParams(params FaceParams, p *runtime.Pinner) C.FaceParams {
	return C.FaceParams{
		fontParams: getCFontParams(params.FontParams, p),
		faceIndex:  C.int(params.FaceIndex),
	}
}

func GetFaceProperties(params FaceParams) (FaceProperties, error) {
	if len(params.FontParams.Buffer) == 0 {
		return FaceProperties{}, fmt.Errorf("invalid font buffer: buffer is empty")
	}

	var p runtime.Pinner
	defer p.Unpin()

	cFaceParams := getCFaceParams(params, &p)

	var cProperties C.FaceProperties

	errCode := C.getFace(cFaceParams, &cProperties)
	if errCode != 0 {
		return FaceProperties{}, fmt.Errorf("%s (error code: 0x%02x)", getErrorMessage(int(errCode)), errCode)
	}

	defer C.freeFace(cProperties)

	return FaceProperties{
		Index:     int(cProperties.index),
		StyleName: C.GoString(cProperties.styleName),
		Monospace: bool(cProperties.monospace),
	}, nil
}

type RenderedCharacter struct {
	Charcode rune

	BitmapBuffer []byte
	BitmapWidth  int
	BitmapHeight int

	LeftShift int
	TopShift  int

	Advance int
}

func goRenderedCharacter(c C.RenderedCharacter) RenderedCharacter {
	var buffer []byte
	if c.bitmapBuffer != nil && c.bitmapWidth > 0 && c.bitmapHeight > 0 {
		size := int(c.bitmapWidth * c.bitmapHeight)
		buffer = C.GoBytes(unsafe.Pointer(c.bitmapBuffer), C.int(size))
	}

	return RenderedCharacter{
		Charcode:     rune(c.charcode),
		BitmapBuffer: buffer,
		BitmapWidth:  int(c.bitmapWidth),
		BitmapHeight: int(c.bitmapHeight),
		LeftShift:    int(c.leftShift),
		TopShift:     int(c.topShift),
		Advance:      int(c.advance),
	}
}

type RenderOutput struct {
	Characters []RenderedCharacter
	Errors     []error

	Monospace bool

	TextHeight int
}

func getCRenderOutput(charactersLength int, p *runtime.Pinner) C.RenderOutput {
	characters := make([]C.RenderedCharacter, charactersLength)
	errors := make([]C.FT_Error, charactersLength)

	if charactersLength > 0 {
		p.Pin(&characters[0])
		p.Pin(&errors[0])
	}

	return C.RenderOutput{
		characters: (*C.RenderedCharacter)(unsafe.Pointer(&characters[0])),
		errors:     (*C.FT_Error)(unsafe.Pointer(&errors[0])),
		length:     0,
	}
}

func goRenderOutput(c C.RenderOutput) RenderOutput {
	var rendered RenderOutput

	rendered.TextHeight = int(c.textHeight)

	cCharactersSlice := unsafe.Slice(c.characters, int(c.length))
	cErrorsSlice := unsafe.Slice(c.errors, int(c.length))

	for i := range c.length {
		cRenderedCharacter, cRenderedErrCode := cCharactersSlice[i], cErrorsSlice[i]

		if cRenderedErrCode != 0 {
			rendered.Errors = append(
				rendered.Errors,
				fmt.Errorf(
					"failed to render character '%c': %s",
					rune(cRenderedCharacter.charcode),
					getErrorMessage(int(cRenderedErrCode)),
				),
			)

			continue
		}

		rendered.Characters = append(
			rendered.Characters,
			goRenderedCharacter(cRenderedCharacter),
		)
	}

	return rendered
}

func Render(face FaceParams, fontSize int, characters []rune) (RenderOutput, error) {
	if len(face.FontParams.Buffer) == 0 {
		return RenderOutput{}, fmt.Errorf("invalid font buffer: buffer is empty")
	}

	var p runtime.Pinner
	defer p.Unpin()

	cFace := getCFaceParams(face, &p)

	cRendered := getCRenderOutput(len(characters), &p)

	errCode := C.render(
		cFace,
		C.int(fontSize),
		(*C.uint)(unsafe.Pointer(&characters[0])),
		C.int(len(characters)),
		&cRendered,
	)
	if errCode != 0 {
		return RenderOutput{}, fmt.Errorf("freetype error: %s (error code: 0x%02x)", getErrorMessage(int(errCode)), errCode)
	}

	defer C.freeRendered(cRendered)

	// @TODO: log rendering errors

	rendered := goRenderOutput(cRendered)
	return rendered, nil
}
