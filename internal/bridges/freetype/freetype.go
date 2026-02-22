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

func newFaceParams(buffer []byte, fontSize int) C.FaceParams {
	faceIndex := 0 // @TODO: multi-face fonts support

	return C.FaceParams{
		buffer:     (*C.uchar)(unsafe.Pointer(&buffer[0])),
		bufferSize: C.int(len(buffer)),
		faceIndex:  C.int(faceIndex),
		fontSize:   C.int(fontSize),
	}
}

func GetFaceProperties(buffer []byte, fontSize int) (*FaceProperties, error) {
	if len(buffer) == 0 {
		return nil, fmt.Errorf("font buffer is empty")
	}
	if fontSize <= 0 {
		return nil, fmt.Errorf("font size invalid: %d", fontSize)
	}

	var p runtime.Pinner
	defer p.Unpin()

	p.Pin(&buffer[0])

	var params C.FaceParams = newFaceParams(buffer, fontSize)
	var properties *C.FaceProperties

	errCode := C.getFaceProperties(&params, &properties)
	if errCode != 0 {
		return nil, fmt.Errorf("%s (error code: 0x%02x)", getErrorMessage(int(errCode)), errCode)
	}

	defer func() {
		if properties != nil {
			C.free(unsafe.Pointer(properties.styleName))
			C.free(unsafe.Pointer(properties.familyName))
			C.free(unsafe.Pointer(properties))
		}
	}()

	return &FaceProperties{
		FamilyName:         C.GoString(properties.familyName),
		StyleName:          C.GoString(properties.styleName),
		Monospace:          bool(properties.monospace),
		MaxCharacterWidth:  int(properties.maxCharacterWidth),
		MaxCharacterHeight: int(properties.maxCharacterHeight),
	}, nil
}
