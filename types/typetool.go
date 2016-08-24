package types

import "github.com/solovev/gopsd/util"

type TypeTool struct {
	Transformation *Matrix
	TextData       *Descriptor
	WarpData       *Descriptor
	// BBox           *RectangleFloat // CHECK TODO Info in spec is not even clear, wtf is 4 * 8? Floats?
}

func ReadTypeTool(reader *util.Reader) *TypeTool {
	tt := new(TypeTool)
	reader.Skip(2) // Version (= 1 for PS 6.0)
	tt.Transformation = ReadMatrix(reader)

	reader.Skip(2) // Text version (= 50 for PS 6.0)
	reader.Skip(4) // Descriptor version (= 16 for PS 6.0)
	tt.TextData = NewDescriptor(reader)

	reader.Skip(2) // Warp version (= 1 for PS 6.0)
	reader.Skip(4) // Descriptor version (= 16 for PS 6.0)
	tt.WarpData = NewDescriptor(reader)

	reader.Skip(32)

	return tt
}

type ObsoleteTypeTool struct {
	Transformation *Matrix
	Faces          []*TypeFace
	Styles         []*TypeStyle

	Type                                   int16
	ScalingFactor, CharacterCount          int32
	HorizontalPlacement, VerticalPlacement int32
	SelectStart, SelectEnd                 int32
	Lines                                  []*TextLine

	ColorSpace int16
	Color      *Color
	AntiAlias  bool
}

type TypeFace struct {
	Mark, Script                                       int16
	FontName, FontFamily, FontStyle                    string
	FontType, DesignVectorAxesCount, DesignVectorValue int32
}

type TypeStyle struct {
	Mark, FaceMark                              int16
	Size, Tracking, Kerning, Leading, BaseShift int32
	AutoKern, Rotate                            bool
}

type TextLine struct {
	CharacterCount                int32
	Orientation, Alignment, Style int16
	ActCharacter                  string
}

func ReadObsoleteTypeTool(reader *util.Reader) *ObsoleteTypeTool {
	tt := new(ObsoleteTypeTool)

	reader.Skip(2) // Version (= 1)
	tt.Transformation = ReadMatrix(reader)

	// Font information
	reader.Skip(2) // Version (= 6)
	tt.Faces = make([]*TypeFace, reader.ReadInt16())
	for i := 0; i < len(tt.Faces); i++ {
		face := new(TypeFace)
		face.Mark = reader.ReadInt16()
		face.FontType = reader.ReadInt32()
		face.FontName = reader.ReadPascalString()
		face.FontFamily = reader.ReadPascalString()
		face.FontStyle = reader.ReadPascalString()
		face.Script = reader.ReadInt16()
		face.DesignVectorAxesCount = reader.ReadInt32()
		face.DesignVectorValue = reader.ReadInt32()
		tt.Faces[i] = face
	}

	// Style information
	tt.Styles = make([]*TypeStyle, reader.ReadInt16())
	for i := 0; i < len(tt.Styles); i++ {
		style := new(TypeStyle)
		style.Mark = reader.ReadInt16()
		style.FaceMark = reader.ReadInt16()
		style.Size = reader.ReadInt32()
		style.Tracking = reader.ReadInt32()
		style.Kerning = reader.ReadInt32()
		style.Leading = reader.ReadInt32()
		style.BaseShift = reader.ReadInt32()
		style.AutoKern = reader.ReadByte() == 1
		reader.Skip(1) // CHECK: Only present in version <= 5
		style.Rotate = reader.ReadByte() == 1
	}

	// Text information
	tt.Type = reader.ReadInt16()
	tt.ScalingFactor = reader.ReadInt32()
	tt.CharacterCount = reader.ReadInt32()
	tt.HorizontalPlacement = reader.ReadInt32()
	tt.VerticalPlacement = reader.ReadInt32()
	tt.SelectStart = reader.ReadInt32()
	tt.SelectEnd = reader.ReadInt32()
	tt.Lines = make([]*TextLine, reader.ReadInt16())
	for i := 0; i < len(tt.Lines); i++ {
		line := new(TextLine)
		line.CharacterCount = reader.ReadInt32()
		line.Orientation = reader.ReadInt16()
		line.Alignment = reader.ReadInt16()
		line.ActCharacter = reader.ReadUnicodeStringLen(1)
		line.Style = reader.ReadInt16()
	}

	// Color information
	tt.ColorSpace = reader.ReadInt16()
	tt.Color = NewRGBAColor(reader)
	tt.AntiAlias = reader.ReadByte() == 1

	return tt
}
