package gopsd

import (
	"bytes"
	"fmt"
	"image"
	"image/jpeg"
)

type ImageResource struct {
	Id   int16
	Name string
	Data interface{}
}

func readResources(doc *Document) {
	length := reader.ReadInt32()

	doc.Resources = make(map[int16]*ImageResource)
	startPos := 0

	for startPos < int(length) {
		ir := new(ImageResource)
		pos := reader.Position

		sign := reader.ReadString32()
		if sign != "8BIM" {
			panic(fmt.Sprintf("Wrong signature of resource №%d!", len(doc.Resources)))
		}

		ir.Id = reader.ReadInt16()
		ir.Name = reader.ReadPascalString()

		size := reader.ReadInt32()
		dataPos := reader.Position
		switch ir.Id {
		case 1033, 1036:
			ir.Data = readResourceThumbnail(size)
		}
		if size%2 != 0 {
			size++
		}
		reader.Skip(int(size) - (reader.Position - dataPos))

		startPos += reader.Position - pos
		doc.Resources[ir.Id] = ir
	}
}

type IRThumbnail struct {
	Width  int32
	Height int32
	Image  image.Image
}

// http://www.adobe.com/devnet-apps/photoshop/fileformatashtml/#50577409_74450
func readResourceThumbnail(size int32) *IRThumbnail {
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
