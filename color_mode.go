package gopsd

import (
	"errors"
	"fmt"
)

type ColorModeDataSection struct {
	Length int32
}

func newColorMode(mode int16) (*ColorModeDataSection, error) {
	c := new(ColorModeDataSection)

	l, err := reader.ReadInt32()
	if err != nil {
		return nil, err
	}
	if mode == 2 {
		if l != 768 {
			return nil, errors.New(fmt.Sprintf("Wrong length of color mode data section: %d! Expected: 768.", l))
		}
		// TODO
		// Indexed color images: length is 768; color data contains the color table for the image, in non-interleaved order.
		err = reader.Skip(l)
		if err != nil {
			return nil, err
		}
	} else if mode == 8 {
		// TODO
		// Duotone images: color data contains the duotone specification (the format of which is not documented).
		// Other applications that read Photoshop files can treat a duotone image as a gray	image, and just preserve the contents of the duotone information when reading and writing the file.
		err = reader.Skip(l)
		if err != nil {
			return nil, err
		}
	} else {
		if l != 0 {
			return nil, errors.New(fmt.Sprintf("Wrong length of color mode data section: %d! Expected: 0.", l))
		}
		err = reader.Skip(l)
		if err != nil {
			return nil, err
		}
	}
	return c, nil
}
