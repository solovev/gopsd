package gopsd

var (
	ColorModes = map[int16]string{
		0: "Bitmap", 1: "Grayscale", 2: "Indexed", 3: "RGB",
		4: "CMYK", 7: "Multichannel", 8: "Duotune", 9: "Lab",
	}
)

type HeaderSection struct {
	Version   int16
	Channels  int16
	Height    int32
	Width     int32
	Depth     int16
	ColorMode int16
}

func newHeader() *HeaderSection {
	header := new(HeaderSection)

	if reader.ReadString32() != "8BPS" {
		panic("Wrong document signature.")
	}

	header.Version = reader.ReadInt16()
	if header.Version != 1 {
		panic("Wrong document version.")
	}

	reader.Skip(6)

	header.Channels = reader.ReadInt16()
	header.Height = reader.ReadInt32()
	header.Width = reader.ReadInt32()
	header.Depth = reader.ReadInt16()

	header.ColorMode = reader.ReadInt16()
	if _, ok := ColorModes[header.ColorMode]; !ok {
		panic("Unknown color mode.")
	}

	return header
}
