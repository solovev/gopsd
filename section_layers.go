package gopsd

import "fmt"

var (
	BlendModeKeys = map[string]string{
		"pass": "Pass through", "norm": "Normal", "diss": "Dissolve",
		"dark": "Darken", "mul": "Multiply", "idiv": "Color burn",
		"lbrn": "Linear burn", "dkCl": "Darker color", "lite": "Lighten",
		"scrn": "Screen", "div": "Color dodge", "lddg": "Linear dodge",
		"lgCl": "Lighter color", "over": "Overlay", "sLit": "Soft light",
		"hLit": "Hard light", "vLit": "Vivid light", "lLit": "Linear light",
		"pLit": "Pin light", "hMix": "Hard mix", "diff": "Difference",
		"smud": "Exclusion", "fsub": "Subtract", "fdiv": "Divide",
		"hue": "Hue", "sat": "Saturation", "colr": "Color", "lum": "Luminosity",
	}
)

func readLayers(doc *Document) {
	var length int64
	if doc.IsLarge {
		length = reader.ReadInt64()
	} else {
		length = int64(reader.ReadInt32())
	}
	if length == 0 {
	}

	var lengthLayers int64
	if doc.IsLarge {
		lengthLayers = reader.ReadInt64()
	} else {
		lengthLayers = int64(reader.ReadInt32())
	}
	//if lengthLayers%2 != 0 {
	//	lengthLayers++
	//}
	lengthLayers = (lengthLayers + 1) & ^0x01

	layerCount := reader.ReadInt16()
	if layerCount < 0 {
		// [TODO] First alpha channel contains the transparency data for the merged result.
		layerCount = -layerCount
	}

	for i := 0; i < int(layerCount); i++ {
		layer := new(Layer)

		layer.Rectangle = reader.ReadRectangle()

		chanCount := reader.ReadInt16()
		for j := 0; j < int(chanCount); j++ {
			channel := new(LayerChannel)

			channel.Id = reader.ReadInt16()
			if doc.IsLarge {
				channel.DataLength = reader.ReadInt64()
			} else {
				channel.DataLength = int64(reader.ReadInt32())
			}
			layer.Channels = append(layer.Channels, channel)
		}

		sign := reader.ReadString32()
		if sign != "8BIM" {
			panic(fmt.Sprintf("Wrong blend mode signature of layer №%d.", i))
		}

		key := reader.ReadString32()
		if mode, ok := BlendModeKeys[key]; ok {
			layer.BlendMode = mode
		}

		layer.Opacity = reader.ReadByte()
		layer.Clipping = reader.ReadByte()
		layer.Flags = reader.ReadByte()
		reader.Skip(1) // Filler

		extraLength := reader.ReadInt32()
		extraPos := reader.Position

		// Mask data
		size := reader.ReadInt32()
		if size != 0 {
			layer.EnclosingMasks = append(layer.EnclosingMasks, reader.ReadRectangle())
			layer.DefaultColor = reader.ReadByte()
			layer.MaskFlags = reader.ReadByte()
			if size == 20 {
				layer.Padding = reader.ReadInt16()
			} else {
				layer.MaskRealFlags = reader.ReadByte()
				layer.MaskBackground = reader.ReadByte()
				layer.EnclosingMasks = append(layer.EnclosingMasks, reader.ReadRectangle())
			}
		}

		// Blending ranges
		blendingLength := reader.ReadInt32()
		layer.BlendingRanges = make([]*LayerBlendingRanges, blendingLength/8)
		for i, value := range layer.BlendingRanges {
			value = new(LayerBlendingRanges)
			if i == 0 {
				value.Name = "Gray"
			} else {
				value.Name = fmt.Sprintf("Channel%d", i-1)
			}
			value.SourceBlack = reader.ReadInt16()
			value.SourceWhite = reader.ReadInt16()
			value.DestBlack = reader.ReadInt16()
			value.DestWhite = reader.ReadInt16()
		}

		layer.Name = reader.ReadPascalString()
		nameLength := len(layer.Name) + 1
		if nameLength%4 != 0 {
			skip := (((nameLength / 4) + 1) * 4) - nameLength
			reader.Skip(skip)
		}

		// Additional information at the end of the layer
		/*for reader.Position < int(extraLength)+extraPos {
			sign = reader.ReadString32()
			if sign != "8BIM" && sign != "8B64" {
				panic(fmt.Sprintf("Wrong additional info signature of layer №%d", i))
			}
			key = reader.ReadString32()
			if doc.IsLarge && stringValueIs(key, "LMsk", "Lr16", "Lr32", "Layr", "Mt16", "Mt32", "Mtrn", "Alph", "FMsk", "lnk2", "FEid", "FXid", "PxSD") {

			}
		}*/
		reader.Skip(int(extraLength) - (reader.Position - extraPos))

		doc.Layers = append(doc.Layers, layer)
	}
}

type Layer struct {
	Rectangle *Rectangle
	Channels  []*LayerChannel
	BlendMode string
	Opacity   byte
	Clipping  byte
	Flags     byte

	// [TODO?] Adjustment layer data
	EnclosingMasks []*Rectangle
	DefaultColor   byte
	MaskFlags      byte
	Padding        int16
	MaskRealFlags  byte
	MaskBackground byte

	// [CHECK] Blending ranges data
	BlendingRanges []*LayerBlendingRanges

	Name string
}

type LayerChannel struct {
	Id         int16
	DataLength int64
}

type LayerBlendingRanges struct {
	Name        string
	SourceBlack int16
	SourceWhite int16
	DestBlack   int16
	DestWhite   int16
}
