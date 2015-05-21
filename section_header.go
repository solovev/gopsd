package gopsd

var (
	ColorModes = map[int16]string{
		0: "Bitmap", 1: "Grayscale", 2: "Indexed", 3: "RGB",
		4: "CMYK", 7: "Multichannel", 8: "Duotune", 9: "Lab",
	}
)

func readHeader(doc *Document) {
	if reader.ReadString32() != "8BPS" {
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
	if !inRange(doc.Channels, 1, 56) {
		panic("The number of channels in the image is out of range.")
	}
	doc.Height = reader.ReadInt32()
	doc.Width = reader.ReadInt32()
	max := 30000
	if doc.IsLarge {
		max *= 10
	}
	if !inRange(doc.Height, 1, max) {
		panic("Document height is out of range.")
	}
	if !inRange(doc.Width, 1, max) {
		panic("Document width is out of range.")
	}

	doc.Depth = reader.ReadInt16()
	if !valueIs(doc.Depth, 1, 8, 16, 32) {
		panic("Wrong value of document depth.")
	}

	cm := reader.ReadInt16()
	if mode, ok := ColorModes[cm]; ok {
		doc.ColorMode = mode
	} else {
		panic("Unknown color mode.")
	}
}
