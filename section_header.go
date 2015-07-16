package gopsd

import "github.com/solovev/gopsd/util"

func readHeader(doc *Document) {
	if reader.ReadString(4) != "8BPS" {
		panic("Wrong document signature.")
	}

	ver := reader.ReadInt16()
	if ver == 2 {
		doc.IsLarge = true
	} else if ver != 1 {
		panic("Wrong document version.")
	}

	reader.Skip(6)

	doc.Channels = reader.ReadInt16()
	if !util.InRange(doc.Channels, 1, 56) {
		panic("The number of channels in the image is out of range.")
	}
	doc.Height = reader.ReadInt32()
	doc.Width = reader.ReadInt32()
	max := 30000
	if doc.IsLarge {
		max *= 10
	}
	if !util.InRange(doc.Height, 1, max) {
		panic("Document height is out of range.")
	}
	if !util.InRange(doc.Width, 1, max) {
		panic("Document width is out of range.")
	}

	doc.Depth = reader.ReadInt16()
	if !util.ValueIs(doc.Depth, 1, 8, 16, 32) {
		panic("Wrong value of document depth.")
	}

	cm := reader.ReadInt16()
	if mode, ok := util.ColorModes[cm]; ok {
		doc.ColorMode = mode
	} else {
		panic("Unknown color mode.")
	}
}
