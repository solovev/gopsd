package gopsd

import (
	"fmt"
)

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
