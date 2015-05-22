package gopsd

import "fmt"

type ImageResource struct {
	UID  int16
	Name string
	Data interface{}
}

func readResources(doc *Document) {
	length := reader.ReadInt32()
	if length == 0 {
		return
	}

	doc.Resources = make(map[int16]*ImageResource)
	pos := 0

	for pos < int(length) {
		ir := new(ImageResource)

		sign := reader.ReadString32()
		if sign != "8BIM" {
			panic(fmt.Sprintf("Wrong signature of resource [%d]!", len(doc.Resources)))
		}
		pos += 4

		id := reader.ReadInt16()
		pos += 2
		ir.UID = id

		name := reader.ReadPascalString()
		pos += len(name) + 1
		ir.Name = name

		size := reader.ReadInt32()
		pos += 4

		switch id {
		default:
			if size%2 != 0 {
				size++
			}
			reader.Skip(size)
		}
		pos += int(size)

		doc.Resources[id] = ir
	}
}
