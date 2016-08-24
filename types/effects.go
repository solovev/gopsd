package types

import "github.com/solovev/gopsd/util"

type ObsoleteEffects struct {
	DropShadow  *ShadowEffect
	InnerShadow *ShadowEffect
	OuterGlow   *GlowEffect
	InnerGlow   *GlowEffect
	Bevel       *BevelEffect
	SolidFill   *SolidFillEffect
}

func ReadObsoleteEffects(reader *util.Reader) *ObsoleteEffects {
	effects := new(ObsoleteEffects)
	reader.Skip(2) // Version (= 0)
	nEffects := int(reader.ReadInt16())
	for i := 0; i < nEffects; i++ {
		reader.Skip(4) // Signature (= "8BIM")
		id := reader.ReadString(4)
		switch id {
		case "cmnS":
			reader.Skip(reader.ReadInt32())
		case "dsdw", "isdw":
			reader.Skip(4)                // Length (41 or 51)
			version := reader.ReadInt32() // Version (0 for PS 5.0 or 2 for 5.5)

			shadow := new(ShadowEffect)
			shadow.Blur = reader.ReadInt32()
			shadow.Intensity = reader.ReadInt32()
			shadow.Angle = reader.ReadInt32()
			shadow.Distance = reader.ReadInt32()
			shadow.ColorSpace = reader.ReadInt16()
			shadow.Color = NewRGBAColor(reader)
			reader.Skip(4) // Blend mode signature
			shadow.BlendMode = reader.ReadString(4)
			shadow.Enabled = reader.ReadByte() == 1
			shadow.SharedEffectAngle = reader.ReadByte() == 1
			shadow.Opacity = reader.ReadByte()
			if version == 2 { // Not stated in spec clearly
				shadow.NativeColorSpace = reader.ReadInt16()
				shadow.NativeColor = NewRGBAColor(reader)
			}
			if id == "dsdw" {
				effects.DropShadow = shadow
			} else {
				effects.InnerShadow = shadow
			}
		case "oglw", "iglw":
			reader.Skip(4) // Length
			version := reader.ReadInt32()

			glow := new(GlowEffect)
			glow.Blur = reader.ReadInt32()
			glow.Intensity = reader.ReadInt32()
			glow.ColorSpace = reader.ReadInt16()
			glow.Color = NewRGBAColor(reader)
			reader.Skip(4) // Blend mode signature
			glow.BlendMode = reader.ReadString(4)
			glow.Enabled = reader.ReadByte() == 1
			glow.Opacity = reader.ReadByte()
			if version == 2 {
				if id == "iglw" {
					glow.Invert = reader.ReadByte() == 1
				}
				glow.NativeColorSpace = reader.ReadInt16()
				glow.NativeColor = NewRGBAColor(reader)
			}
			if id == "oglw" {
				effects.OuterGlow = glow
			} else {
				effects.InnerGlow = glow
			}
		case "bevl":
			reader.Skip(4) // Length
			version := reader.ReadInt32()

			bevel := new(BevelEffect)
			bevel.Angle = reader.ReadInt32()
			bevel.Strength = reader.ReadInt32()
			bevel.Blur = reader.ReadInt32()
			reader.Skip(4) // Blend mode signature
			bevel.HighlightBlendMode = reader.ReadString(4)
			reader.Skip(4) // Blend mode signature
			bevel.ShadowBlendMode = reader.ReadString(4)
			bevel.HighlightColorSpace = reader.ReadInt16()
			bevel.HighlightColor = NewRGBAColor(reader)
			bevel.ShadowColorSpace = reader.ReadInt16()
			bevel.ShadowColor = NewRGBAColor(reader)
			bevel.BevelStyle = reader.ReadByte()
			bevel.HighlightOpacity = reader.ReadByte()
			bevel.ShadowOpacity = reader.ReadByte()
			bevel.Enabled = reader.ReadByte() == 1
			bevel.SharedEffectAngle = reader.ReadByte() == 1
			bevel.Up = reader.ReadByte() == 1
			if version == 2 {
				bevel.RealHighlightColorSpace = reader.ReadInt16()
				bevel.RealHighlightColor = NewRGBAColor(reader)
				bevel.RealShadowColorSpace = reader.ReadInt16()
				bevel.RealShadowColor = NewRGBAColor(reader)
			}
			effects.Bevel = bevel
		case "sofi":
			reader.Skip(4) // Length
			reader.Skip(4) // Version (= 2)

			fill := new(SolidFillEffect)
			fill.BlendMode = reader.ReadString(4)
			fill.ColorSpace = reader.ReadInt16()
			fill.Color = NewRGBAColor(reader)
			fill.Opacity = reader.ReadByte()
			fill.Enabled = reader.ReadByte() == 1
			fill.NativeColorSpace = reader.ReadInt16()
			fill.NativeColor = NewRGBAColor(reader)

			effects.SolidFill = fill
		}
	}
	return effects
}

type ShadowEffect struct {
	Blur, Intensity, Angle, Distance int32
	ColorSpace, NativeColorSpace     int16
	Color, NativeColor               *Color
	BlendMode                        string
	Enabled, SharedEffectAngle       bool
	Opacity                          byte
}

type GlowEffect struct {
	Blur, Intensity              int32
	ColorSpace, NativeColorSpace int16
	Color, NativeColor           *Color
	BlendMode                    string
	Enabled, Invert              bool
	Opacity                      byte
}

type BevelEffect struct {
	Angle, Strength, Blur                       int32
	HighlightBlendMode, ShadowBlendMode         string
	HighlightColorSpace, ShadowColorSpace       int16
	HighlightColor, ShadowColor                 *Color
	BevelStyle, HighlightOpacity, ShadowOpacity byte
	Enabled, SharedEffectAngle, Up              bool

	RealHighlightColorSpace, RealShadowColorSpace int16
	RealHighlightColor, RealShadowColor           *Color
}

type SolidFillEffect struct {
	BlendMode                    string
	ColorSpace, NativeColorSpace int16
	Color, NativeColor           *Color
	Opacity                      byte
	Enabled                      bool
}
