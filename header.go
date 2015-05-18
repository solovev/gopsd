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
		return nil, errors.New(fmt.Sprintf("Wrong document version: %s!", ver))
	}
	header.Version = ver

	return header, nil
}
