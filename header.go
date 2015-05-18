package gopsd

import (
	"errors"
	"fmt"
)

type HeaderSection struct {
	Signature string
	Version   int16
	Channels  int16
	Height    int32
	Width     int32
	Depth     int16
	ColorMode int16
}

func newHeader() (*HeaderSection, error) {
	header := new(HeaderSection)

	sign, err := reader.ReadString32()
	if err != nil {
		return nil, err
	}
	if sign != "8BPS" {
		return nil, errors.New(fmt.Sprintf("Wrong document signature: %s!", sign))
	}
	header.Signature = sign

	ver, err := reader.ReadInt16()
	if err != nil {
		return nil, err
	}
	if ver != 1 {
		return nil, errors.New(fmt.Sprintf("Wrong document version: %d!", ver))
	}
	header.Version = ver

	err = reader.Skip(6)
	if err != nil {
		return nil, err
	}

	ch, err := reader.ReadInt16()
	if err != nil {
		return nil, err
	}
	if ch < 1 || ch > 56 {
		return nil, errors.New(fmt.Sprintf("Wrong number of channels: %d! Supported range is 1 to 56.", ch))
	}
	header.Channels = ch

	h, err := reader.ReadInt32()
	if err != nil {
		return nil, err
	}
	if h < 1 || h > 30000 {
		return nil, errors.New(fmt.Sprintf("Wrong document height: %d! Supported range is 1 to 30000.", h))
	}
	header.Height = h

	w, err := reader.ReadInt32()
	if err != nil {
		return nil, err
	}
	if w < 1 || w > 30000 {
		return nil, errors.New(fmt.Sprintf("Wrong document width: %d! Supported range is 1 to 30000.", w))
	}
	header.Width = w

	d, err := reader.ReadInt16()
	if err != nil {
		return nil, err
	}
	if d != 1 && d != 8 && d != 16 && d != 32 {
		return nil, errors.New(fmt.Sprintf("Wrong document depth: %d! Supported values are 1, 8, 16 and 32.", d))
	}
	header.Depth = d

	c, err := reader.ReadInt16()
	if err != nil {
		return nil, err
	}
	if c < 0 || c > 9 {
		return nil, errors.New(fmt.Sprintf("Wrong color mode: %d! Supported values are: Bitmap = 0; Grayscale = 1; Indexed = 2; RGB = 3; CMYK = 4; Multichannel = 7; Duotone = 8; Lab = 9.", c))
	}
	header.ColorMode = c

	return header, nil
}
