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

type FaceProperties struct {
	FamilyName string
	StyleName  string

	Monospace bool

	MaxCharacterWidth  int
	MaxCharacterHeight int
}

type FaceParams struct {
	FontBuffer []byte
	FontSize   int
	FaceIndex  int
}

func newFaceParams(params FaceParams) (*C.FaceParams, error) {
	if len(params.FontBuffer) == 0 {
		return nil, fmt.Errorf("font buffer is empty")
	}
	if params.FontSize <= 0 {
		return nil, fmt.Errorf("font size invalid: %d", params.FontSize)
	}

	return &C.FaceParams{
		buffer:     (*C.uchar)(unsafe.Pointer(&params.FontBuffer[0])),
		bufferSize: C.int(len(params.FontBuffer)),
		faceIndex:  C.int(params.FaceIndex),
		fontSize:   C.int(params.FontSize),
	}, nil
}

func GetFaceProperties(params FaceParams) (*FaceProperties, error) {
	var p runtime.Pinner
	p.Pin(&params.FontBuffer[0])
	defer p.Unpin()

	cParams, error := newFaceParams(params)
	if error != nil {
		return nil, error
	}

	var cProperties C.FaceProperties

	errCode := C.getFaceProperties(cParams, &cProperties)
	if errCode != 0 {
		return nil, fmt.Errorf("%s (error code: 0x%02x)", getErrorMessage(int(errCode)), errCode)
	}

	defer C.freeFaceProperties(&cProperties)

	return &FaceProperties{
		FamilyName:         C.GoString(cProperties.familyName),
		StyleName:          C.GoString(cProperties.styleName),
		Monospace:          bool(cProperties.monospace),
		MaxCharacterWidth:  int(cProperties.maxCharacterWidth),
		MaxCharacterHeight: int(cProperties.maxCharacterHeight),
	}, nil
}

type RenderedCharacter struct {
	BitmapBuffer []byte
	BitmapWidth  int
	BitmapHeight int

	LeftShift int
	TopShift  int

	Advance int
}

func RenderCharacters(params FaceParams, characters []rune) ([]*RenderedCharacter, error) {
	var p runtime.Pinner
	p.Pin(&params.FontBuffer[0])
	defer p.Unpin()

	cParams, error := newFaceParams(params)
	if error != nil {
		return nil, error
	}

	length := len(characters)
	if length == 0 {
		return nil, nil
	}

	cCharacters := (*C.uint)(unsafe.Pointer(&characters[0]))
	cLength := C.int(length)

	cRenderedChars := make([]*C.RenderedCharacter, length)
	cRenderingErrors := make([]C.FT_Error, length)

	errCode := C.renderCharacters(cParams, cCharacters, cLength, &cRenderedChars[0], &cRenderingErrors[0])
	if errCode != 0 {
		return nil, fmt.Errorf("freetype error: %s (error code: 0x%02x)", getErrorMessage(int(errCode)), errCode)
	}

	defer C.freeRenderedCharacters(&cRenderedChars[0], cLength)

	// @TODO: log rendering errors

	renderedCharacters := make([]*RenderedCharacter, length)

	for i := range length {
		cRC := cRenderedChars[i]
		if cRC == nil {
			renderedCharacters[i] = nil
			continue
		}

		var buffer []byte
		if cRC.bitmapBuffer != nil && cRC.bitmapWidth > 0 && cRC.bitmapHeight > 0 {
			size := int(cRC.bitmapWidth * cRC.bitmapHeight)
			buffer = C.GoBytes(unsafe.Pointer(cRC.bitmapBuffer), C.int(size))
		}

		renderedCharacters[i] = &RenderedCharacter{
			BitmapBuffer: buffer,
			BitmapWidth:  int(cRC.bitmapWidth),
			BitmapHeight: int(cRC.bitmapHeight),
			LeftShift:    int(cRC.leftShift),
			TopShift:     int(cRC.topShift),
			Advance:      int(cRC.advance),
		}
	}

	return renderedCharacters, nil
}
