package font

import (
	"fmt"
	"os"

	"github.com/ddmytro-m/asciitor/internal/bridges/freetype"
)

type Font struct {
	fontBuffer []byte

	familyName  string
	facesAmount int

	loaded bool
}

func NewFont(buffer []byte) (*Font, error) {
	font := Font{fontBuffer: buffer}

	properties, err := freetype.GetFontProperties(font.GetParams())

	if err != nil {
		return nil, err
	}

	font.familyName = properties.FamilyName
	font.facesAmount = properties.FacesAmount

	font.loaded = true

	return &font, nil
}

func NewFontFromFile(file string) (*Font, error) {
	data, err := os.ReadFile(file)
	if err != nil {
		return nil, fmt.Errorf("failed to open font file: \"%s\"", file)
	}

	return NewFont(data)
}

func (f *Font) GetParams() freetype.FontParams {
	return freetype.FontParams{
		Buffer: f.fontBuffer,
	}
}

func (f *Font) GetFace(index int) (*Face, error) {
	faceProperties, err := freetype.GetFaceProperties(freetype.FaceParams{FontParams: f.GetParams(), FaceIndex: index})
	if err != nil {
		return nil, err
	}

	face := newFace(f, faceProperties)
	return face, nil
}

func (f *Font) GetFaces() ([]Face, error) {
	faces, err := freetype.GetFaces(f.GetParams())
	if err != nil {
		return nil, err
	}

	outFaces := make([]Face, len(faces))
	for i, properties := range faces {
		outFaces[i] = *newFace(f, properties)
	}

	return outFaces, nil
}
