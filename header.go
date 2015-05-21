package gopsd

var (
	LargeDocument bool

	ColorModes = map[int16]string{
		0: "Bitmap", 1: "Grayscale", 2: "Indexed", 3: "RGB",
		4: "CMYK", 7: "Multichannel", 8: "Duotune", 9: "Lab",
	}
)

type HeaderSection struct {
	Channels  int16
	Height    int32
	Width     int32
	Depth     int16
	ColorMode string
}

func newHeader() *HeaderSection {
	header := new(HeaderSection)

	if reader.ReadString32() != "8BPS" {
		panic("Wrong document signature.")
	}

	ver := reader.ReadInt16()
	if ver == 2 {
		LargeDocument = true
	} else if ver != 1 {
		panic("Wrong document version.")
	}

	reader.Skip(6)

	header.Channels = reader.ReadInt16()
	if !inRange(header.Channels, 1, 56) {
		panic("The number of channels in the image is out of range.")
	}
	header.Height = reader.ReadInt32()
	header.Width = reader.ReadInt32()
	max := 30000
	if LargeDocument {
		max *= 10
	}
	if !inRange(header.Height, 1, max) {
		panic("Document height is out of range.")
	}
	if !inRange(header.Width, 1, max) {
		panic("Document width is out of range.")
	}

	header.Depth = reader.ReadInt16()
	if !valueIs(header.Depth, 1, 8, 16, 32) {
		panic("Wrong value of header depth.")
	}

	cm := reader.ReadInt16()
	if mode, ok := ColorModes[cm]; ok {
		header.ColorMode = mode
	} else {
		panic("Unknown color mode.")
	}

	return header
}
