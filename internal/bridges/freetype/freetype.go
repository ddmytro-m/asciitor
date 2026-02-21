package freetype

//#cgo pkg-config: freetype2
//#include "freetype.h"
import "C"
import (
	"fmt"
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

func GetFaceProperties(buffer []byte, fontSize int) (*FaceProperties, error) {
	faceIndex := 0 // @TODO: multi-face fonts support

	var properties *C.FaceProperties

	err := C.getFaceProperties((*C.uchar)(unsafe.Pointer(&buffer[0])), C.int(len(buffer)), C.int(faceIndex), C.int(fontSize), &properties)
	if err != 0 {
		return nil, fmt.Errorf("%s (error code: 0x%02x)", getErrorMessage(int(err)), err)
	}

	defer func() {
		C.free(unsafe.Pointer(properties.styleName))
		C.free(unsafe.Pointer(properties.familyName))
		C.free(unsafe.Pointer(properties))
	}()

	return &FaceProperties{
		FamilyName:         C.GoString(properties.familyName),
		StyleName:          C.GoString(properties.styleName),
		Monospace:          bool(properties.monospace),
		MaxCharacterWidth:  int(properties.maxCharacterWidth),
		MaxCharacterHeight: int(properties.maxCharacterHeight),
	}, nil
}
