package gopsd

import (
	"errors"
	"fmt"
)

var (
	BlendModeKeys = map[string]string{
		"pass": "Pass through", "norm": "Normal", "diss": "Dissolve",
		"dark": "Darken", "mul": "Multiply", "idiv": "Color burn",
		"lbrn": "Linear burn", "dkCl": "Darker color", "lite": "Lighten",
		"scrn": "Screen", "div": "Color dodge", "lddg": "Linear dodge",
		"lgCl": "Lighter color", "over": "Overlay", "sLit": "Soft light",
		"hLit": "Hard light", "vLit": "Vivid light", "lLit": "Linear light",
		"pLit": "Pin light", "hMix": "Hard mix", "diff": "Difference",
		"smud": "Exclusion", "fsub": "Subtract", "fdiv": "Divide",
		"hue": "Hue", "sat": "Saturation", "colr": "Color", "lum": "Luminosity",
	}
)

type LayerAndMaskInfoSection struct {
	Layers      []*LayerRecord
	MergedAlpha bool
}

type LayerRecord struct {
	Rectangle    *Rectangle
	Channels     []*LayerChannel
	BlendModeKey string
	Opacity      byte
	Clipping     byte
	Flags        byte
	// Layer mask / adjustment layer data
	Padding        int16
	MaskRectangle  *Rectangle
	DefaultColor   byte
	MaskFlags      byte
	MaskRealFlags  byte
	MaskBackground byte
	MaskRectangle2 *Rectangle
}

func (l *LayerRecord) readChannels() error {
	// Reading number of channels in the layer
	chN, err := reader.ReadInt16()
	if err != nil {
		return err
	}

	// Reading information about each channel
	for j := 0; j < int(chN); j++ {
		ch := new(LayerChannel)

		// Reading channel ID
		chId, err := reader.ReadInt16()
		if err != nil {
			return err
		}
		ch.ChannelID = chId

		// Reading length of channel data
		chL, err := reader.ReadInt32()
		if err != nil {
			return err
		}
		ch.Length = chL

		l.Channels = append(l.Channels, ch)
	}

	return nil
}

func (l *LayerRecord) readExtraData() error {
	// Reading length of the extra data
	xtrl, err := reader.ReadInt32()
	if err != nil {
		return err
	}

	// Reading size of the adjustment data
	size, err := reader.ReadInt32()
	if err != nil {
		return err
	}
	if size > 0 {
		// Reading enclosing layer mask
		r, err := reader.ReadRectangle()
		if err != nil {
			return err
		}
		l.MaskRectangle = r

		// Reading default color
		clr, err := reader.ReadByte()
		if err != nil {
			return err
		}
		l.DefaultColor = clr

		// [TODO] Mask parameters (flags bit 4)
		// Reading flags
		fl, err := reader.ReadByte()
		if err != nil {
			return err
		}
		l.MaskFlags = fl

		if xtrl == 20 {
			pad, err := reader.ReadInt16()
			if err != nil {
				return err
			}
		}
	}

	return nil
}

type LayerChannel struct {
	ChannelID int16
	Length    int32
}

func newLayerAndMaskInfo() (*LayerAndMaskInfoSection, error) {
	lam := new(LayerAndMaskInfoSection)

	// Reading length of the layer and mask information section
	sl, err := reader.ReadInt32()
	if err != nil {
		return nil, err
	}

	// Reading length of the layers info section, rounded up to a multiple of 2
	l, err := reader.ReadInt32()
	if err != nil {
		return nil, err
	}

	// Reading layer count
	n, err := reader.ReadInt16()
	if err != nil {
		return nil, err
	}
	//  If it is a negative number, its absolute value is the number of layers and the first alpha channel contains the transparency data for the merged result.
	if n < 0 {
		lam.MergedAlpha = true
		n = -n
	}

	// Reading information about each layer
	for i := 0; i < int(n); i++ {
		layer := new(LayerRecord)

		// Reading layer's rectangle
		r, err := reader.ReadRectangle()
		if err != nil {
			return nil, err
		}
		layer.Rectangle = r

		// Reading layer's channels
		if err = layer.readChannels(); err != nil {
			return nil, err
		}

		// Reading blend mode signature
		sign, err := reader.ReadString32()
		if err != nil {
			return nil, err
		}
		if sign != "8BIM" {
			return nil, errors.New(fmt.Sprintf("Wrong blend mode signature: %s!", sign))
		}
		fmt.Println(sign)

		// Reading blend mode key
		key, err := reader.ReadString32()
		if err != nil {
			return nil, err
		}
		if _, ok := BlendModeKeys[key]; !ok {
			return nil, errors.New(fmt.Sprintf("Unknown blend mode key: %s!", key))
		}
		layer.BlendModeKey = key

		// Reading opacity
		op, err := reader.ReadByte()
		if err != nil {
			return nil, err
		}
		layer.Opacity = op

		// Reading clipping
		cl, err := reader.ReadByte()
		if err != nil {
			return nil, err
		}
		layer.Clipping = cl

		// Reading flags
		fl, err := reader.ReadByte()
		if err != nil {
			return nil, err
		}
		layer.Flags = fl

		// Reading filler
		_, err = reader.ReadByte()
		if err != nil {
			return nil, err
		}

		lam.Layers = append(lam.Layers, layer)
	}

	return lam, nil
}
