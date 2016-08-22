package gopsd

import (
	"bytes"
	"fmt"
	"image"
	"image/jpeg"

	"github.com/solovev/gopsd/types"
	"github.com/solovev/gopsd/util"
)

type IRThumbnail struct {
	Width  int32
	Height int32
	Image  image.Image
}

type IRPrintStyle struct {
	DescriptorVersion int32
	Descriptor        *types.Descriptor
}

type IRAspectRatio struct {
	Version int32
	Ratio   float64
}

// http://www.adobe.com/devnet-apps/photoshop/fileformatashtml/#50577409_74450
func ReadResourceThumbnail(reader *util.Reader) *IRThumbnail {
	thumb := new(IRThumbnail)

	format := reader.ReadInt32()
	thumb.Width = reader.ReadInt32()
	thumb.Height = reader.ReadInt32()

	reader.ReadInt32() // Widthbytes
	reader.ReadInt32() // Total size
	comprSize := reader.ReadInt32()

	reader.ReadInt16() // Bits per pixel
	reader.ReadInt16() // Number of planes

	switch format {
	case 0:
	case 1:
		img, err := jpeg.Decode(bytes.NewReader(reader.ReadBytes(comprSize)))
		if err != nil {
			panic(err)
		}
		thumb.Image = img
	default:
	}

	return thumb
}

func ReadResourcePrintStyle(reader *util.Reader) *IRPrintStyle {
	style := new(IRPrintStyle)

	style.DescriptorVersion = reader.ReadInt32()
	if style.DescriptorVersion == 16 {
		style.Descriptor = types.NewDescriptor(reader)
	}

	return style
}

// TODO
func ReadResourceAspectRatio(reader *util.Reader) *IRAspectRatio {
	ratio := new(IRAspectRatio)

	ratio.Version = reader.ReadInt32()
	//ratio.Ratio = reader.ReadFloat64()

	return ratio
}

func readResources(doc *Document) {
	length := reader.ReadInt32()

	doc.Resources = make(map[int16]interface{})
	startPos := 0

	for startPos < int(length) {
		pos := reader.Position

		sign := reader.ReadString(4)
		if sign != "8BIM" {
			panic(fmt.Sprintf("Wrong signature of resource #%d!", len(doc.Resources)))
		}

		id := reader.ReadInt16()
		// Resource name [CHECK]
		reader.ReadPascalString()

		size := reader.ReadInt32()
		dataPos := reader.Position

		switch id {
		case 1033, 1036:
			doc.Resources[id] = ReadResourceThumbnail(reader)
		case 1083:
			doc.Resources[id] = ReadResourcePrintStyle(reader)
		case 1064:
			doc.Resources[id] = ReadResourceAspectRatio(reader)
		default:
			doc.Resources[id] = nil
		}
		if size%2 != 0 {
			size++
		}
		reader.Skip(int(size) - (reader.Position - dataPos))

		startPos += reader.Position - pos
	}
}
