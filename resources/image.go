package resources

import (
	"bytes"
	"image"
	"image/jpeg"

	"github.com/solovev/gopsd/util"
)

type IRThumbnail struct {
	Width  int32
	Height int32
	Image  image.Image
}

type IRPrintStyle struct {
	DescriptorVersion int32
	Descriptor        *util.Descriptor
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
		style.Descriptor = util.NewDescriptor(reader)
	}

	return style
}
