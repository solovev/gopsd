package gopsd

import (
	"fmt"
	"image"
	"image/color"
	"math"

	"github.com/solovev/gopsd/types"
	"github.com/solovev/gopsd/util"
)

func readLayers(doc *Document) {
	var length int64
	if doc.IsLarge {
		length = reader.ReadInt64()
	} else {
		length = int64(reader.ReadInt32())
	}
	pos := reader.Position

	var lengthLayers int64
	if doc.IsLarge {
		lengthLayers = reader.ReadInt64()
	} else {
		lengthLayers = int64(reader.ReadInt32())
	}
	lengthLayers = lengthLayers + 1 & ^0x01

	layerCount := reader.ReadInt16()
	if layerCount < 0 {
		// [TODO] First alpha channel contains the transparency data for the merged result.
		layerCount = -layerCount
	}

	for i := 0; i < int(layerCount); i++ {
		layer := new(Layer)
		layer.Type = TypeUnspecified
		layer.Rectangle = types.NewRectangle(reader)

		chanCount := reader.ReadInt16()
		for j := 0; j < int(chanCount); j++ {
			channel := new(LayerChannel)

			channel.ID = reader.ReadInt16()
			if doc.IsLarge {
				channel.Length = reader.ReadInt64()
			} else {
				channel.Length = int64(reader.ReadInt32())
			}
			layer.Channels = append(layer.Channels, channel)
		}

		sign := reader.ReadString(4)
		if sign != "8BIM" {
			panic(fmt.Sprintf("Wrong blend mode signature of layer [#%d].", i))
		}

		key := reader.ReadString(4)
		if mode, ok := util.BlendModeKeys[key]; ok {
			layer.BlendMode = mode
		}

		layer.Opacity = byte(math.Ceil(float64(reader.ReadByte()) / 255 * 100))
		layer.Clipping = reader.ReadByte()
		layer.Flags = reader.ReadByte()
		reader.Skip(1) // Filler

		extraLength := reader.ReadInt32()
		extraPos := reader.Position

		// Mask data
		size := reader.ReadInt32()
		if size != 0 {
			layer.EnclosingMasks = append(layer.EnclosingMasks, types.NewRectangle(reader))
			layer.DefaultColor = reader.ReadByte()
			layer.MaskFlags = reader.ReadByte()
			if size == 20 {
				layer.Padding = reader.ReadInt16()
			} else {
				layer.MaskRealFlags = reader.ReadByte()
				layer.MaskBackground = reader.ReadByte()
				layer.EnclosingMasks = append(layer.EnclosingMasks, types.NewRectangle(reader))
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

		// Name. Pascal string, padded to a multiple of 4 bytes
		layer.Name = reader.ReadPascalString()
		nameLength := len(layer.Name) + 1
		if nameLength%4 != 0 {
			skip := 4 - nameLength%4
			reader.Skip(skip)
		}

		// Additional information at the end of the layer
		index := 0
		for reader.Position < int(extraLength)+extraPos {
			sign = reader.ReadString(4)
			if sign != "8BIM" && sign != "8B64" {
				panic(fmt.Sprintf("[Layer: %s] Wrong signature of additional info [#%d]", layer.Name, index))
			}
			key = reader.ReadString(4)
			layer.DataKeys = append(layer.DataKeys, key)

			var dataLength int64
			if doc.IsLarge && util.StringValueIs(key, "LMsk", "Lr16", "Lr32", "Layr", "Mt16", "Mt32", "Mtrn", "Alph", "FMsk", "lnk2", "FEid", "FXid", "PxSD") {
				dataLength = reader.ReadInt64()
			} else {
				dataLength = int64(reader.ReadInt32())
			}
			dataLength = dataLength + 1 & ^0x01
			dataPos := reader.Position

			switch key {
			case "tySh":
				layer.ObsoleteTypeTool = types.ReadObsoleteTypeTool(reader)
			case "TySh":
				layer.TypeTool = types.ReadTypeTool(reader)
			case "luni":
				layer.Name = reader.ReadUnicodeString()
			case "lnsr": // layr / bgnd
				switch reader.ReadString(4) {
				case "layr":
					layer.Type = TypeDefault
				case "shap":
					layer.Type = TypeShape
				case "bgnd":
					layer.Type = TypeBackground
				case "rend":
					layer.Type = TypeRenderObject
				case "lset":
					layer.Type = TypeLockSet
				}
			case "lyid":
				layer.ID = reader.ReadInt32()
			case "clbl":
				layer.BlendClippedElements = reader.ReadByte() == 1
				reader.Skip(3)
			case "infx":
				layer.BlendInteriorElements = reader.ReadByte() == 1
				reader.Skip(3)
			case "knko":
				layer.Knockout = reader.ReadByte() == 1
				reader.Skip(3)
			case "lspf":
				layer.ProtectionFlags = reader.ReadInt32()
			case "lclr":
				layer.SheetColor = types.NewRGBAColor(reader)
			case "fxrp":
				point := make([]float64, 2)
				point[0] = reader.ReadFloat64()
				point[1] = reader.ReadFloat64()
				layer.ReferencePoint = point
			case "lsct":
				sectionType := reader.ReadInt32()
				if sectionType == 3 {
					layer.IsSectionDivider = true
				} else if sectionType > 0 {
					layer.IsFolder = true
				}
				if dataLength >= 12 {
					if reader.ReadString(4) != "8BIM" {
						panic(fmt.Sprintf("Wrong section signature of layer %s", layer.Name))
					}
					key := reader.ReadString(4)
					if mode, ok := util.BlendModeKeys[key]; ok {
						layer.BlendMode = mode // Overriding (as group)
					}
				}
				if dataLength >= 16 {
					layer.IsSceneGroup = reader.ReadInt32() > 0
				}
			case "lsdk": // Inserted layer group (not present in spec)
				sectionType := reader.ReadInt32()
				if sectionType == 3 {
					layer.IsSectionDivider = true
				} else if sectionType > 0 {
					layer.IsFolder = true
				}
			case "lfx2":
				reader.ReadInt32()
				reader.ReadInt32()
				types.NewDescriptor(reader)
			case "vogk": // TODO (Shape bounding box)
				reader.Skip(4) // Version (= 1 for PS CC)
				reader.Skip(4) // Version (= 16)
				layer.VectorOriginData = types.NewDescriptor(reader)
			case "vmsk", "vsms":
				reader.Skip(4) // Version (= 3 for PS 6.0)
				flags := uint32(reader.ReadInt32())
				vectorMask := new(LayerVectorMask)
				vectorMask.IsInverted = (flags & (1 << 0)) > 0
				vectorMask.IsLinked = (flags & (1 << 1)) == 0
				vectorMask.IsDisabled = (flags & (1 << 2)) > 0
				vectorMask.Path = types.ReadPath(doc.Width, doc.Height, reader.ReadBytes(dataLength-8))
				layer.VectorMask = vectorMask
			default:
				reader.Skip(dataLength)
			}
			reader.Skip(dataPos + int(dataLength) - reader.Position)
			index++
		}
		// [CHECK] Not needed
		reader.Skip(int(extraLength) - (reader.Position - extraPos))
		doc.Layers = append(doc.Layers, layer)
	}

	for _, layer := range doc.Layers {
		width := int(layer.Rectangle.Width)
		height := int(layer.Rectangle.Height)

		for _, channel := range layer.Channels {
			compression := reader.ReadInt16()
			switch compression {
			case 0:
				channel.Data = reader.ReadSignedBytes(width * height)
			case 1:
				var result []int8
				scanLines := make([]int16, height)
				for i := range scanLines {
					scanLines[i] = reader.ReadInt16()
				}
				for i := range scanLines {
					line := util.UnpackRLEBits(reader.ReadSignedBytes(scanLines[i]), width)
					result = append(result, line...)
				}
				channel.Data = result
			default:
				panic(fmt.Sprintf("[Layer: %s] Unknown compression method of channel [id: %d]", layer.Name, channel.ID))
			}
		}
	}
	reader.Skip(int(length) - (reader.Position - pos))
}

func (l Layer) ToString() string {
	return fmt.Sprintf("%s: %s", l.Name, l.Rectangle.ToString())
}

type Layer struct {
	ID        int32
	Name      string
	Rectangle *types.Rectangle
	Channels  []*LayerChannel `json:"-"`
	BlendMode string          `json:"-"`
	Opacity   byte            `json:"-"`
	Clipping  byte            `json:"-"`
	Flags     byte            `json:"-"`

	// [TODO?] Adjustment layer data
	EnclosingMasks []*types.Rectangle `json:"-"`
	DefaultColor   byte               `json:"-"`
	MaskFlags      byte               `json:"-"`
	Padding        int16              `json:"-"`
	MaskRealFlags  byte               `json:"-"`
	MaskBackground byte               `json:"-"`

	// [CHECK] Blending ranges data, empty name
	BlendingRanges []*LayerBlendingRanges `json:"-"`

	Type                  LayerType    `json:"-"`
	BlendClippedElements  bool         `json:"-"`
	BlendInteriorElements bool         `json:"-"`
	Knockout              bool         `json:"-"`
	ProtectionFlags       int32        `json:"-"`
	SheetColor            *types.Color `json:"-"`
	ReferencePoint        []float64    `json:"-"`
	IsSceneGroup          bool         `json:"-"`
	IsFolder              bool         `json:"-"`
	IsSectionDivider      bool         `json:"-"`
	DataKeys              []string

	VectorMask       *LayerVectorMask  `json:"-"`
	VectorOriginData *types.Descriptor `json:"-"`

	ObsoleteTypeTool *types.ObsoleteTypeTool `json:"-"`
	TypeTool         *types.TypeTool         `json:"-"`

	Parent   *Layer
	Children []*Layer
}

func (l *Layer) IsText() bool {
	return l.ObsoleteTypeTool != nil || l.TypeTool != nil
}

func (l *Layer) GetImage() (image.Image, error) {
	width := int(l.Rectangle.Width)
	height := int(l.Rectangle.Height)

	if width == 0 || height == 0 {
		return nil, nil
	}

	image := image.NewRGBA(image.Rect(0, 0, width, height))
	switch len(l.Channels) {
	case 3: // RGB
		// [TODO]
	case 4, 5:
		c := l.Channels
		for x := 0; x < width; x++ {
			for y := 0; y < height; y++ {
				i := x + (y * width)
				red := byte(c[1].Data[i])
				green := byte(c[2].Data[i])
				blue := byte(c[3].Data[i])
				alpha := byte(c[0].Data[i])
				image.Set(x, y, color.RGBA{red, green, blue, alpha})
			}
		}
	}
	return image, nil
}

type LayerVectorMask struct {
	IsInverted, IsLinked, IsDisabled bool
	Path                             *types.Path
}

// LayerChannel stores color data of channel.
// Channel IDs:
//		0 = red, 1 = green, 2 = blue;
//		-1 = transparency mask
//		-2 = user supplied layer mask
//		-3 real user supplied layer mask
type LayerChannel struct {
	ID int16
	// [CHECK]
	Length int64
	Data   []int8
}

type LayerBlendingRanges struct {
	Name        string
	SourceBlack int16
	SourceWhite int16
	DestBlack   int16
	DestWhite   int16
}

// [TODO] Not impl yet
type GlobalLayerMask struct {
	OverlayColorSpace int16
	ColorComponents   []int16
	Opacity           int16
	Kind              byte
}

type LayerType int

const (
	TypeDefault LayerType = iota
	TypeShape
	TypeBackground
	TypeRenderObject
	TypeLockSet
	TypeUnspecified
)
