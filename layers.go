package gopsd

import (
	"errors"
	"fmt"
)

type LayerAndMaskInfoSection struct {
	Length       int32
	LayersLength int32
	Layers       []*LayerRecord
}

type LayerRecord struct {
	Rectangle *Rectangle
	Channels  []*LayerChannel
}

type LayerChannel struct {
	ChannelID int16
	Length    int32
}

func newLayerAndMaskInfo() (*LayerAndMaskInfoSection, error) {
	lam := new(LayerAndMaskInfoSection)

	sl, err := reader.ReadInt32()
	if err != nil {
		return nil, err
	}
	lam.Length = sl

	l, err := reader.ReadInt32()
	if err != nil {
		return nil, err
	}
	lam.LayersLength = l

	n, err := reader.ReadInt16()
	if err != nil {
		return nil, err
	}

	for i := 0; i < int(n); i++ {
		layer := new(LayerRecord)

		r, err := reader.ReadRectangle()
		if err != nil {
			return nil, err
		}
		layer.Rectangle = r

		chN, err := reader.ReadInt16()
		if err != nil {
			return nil, err
		}

		for j := 0; j < int(chN); j++ {
			ch := new(LayerChannel)

			chId, err := reader.ReadInt16()
			if err != nil {
				return nil, err
			}
			if chId < -3 || chId > 4 {
				return nil, errors.New(fmt.Sprintf("Wrong channel id: %d!", chId))
			}
			ch.ChannelID = chId

			chL, err := reader.ReadInt32()
			if err != nil {
				return nil, err
			}
			ch.Length = chL

			layer.Channels = append(layer.Channels, ch)
		}

		lam.Layers = append(lam.Layers, layer)
	}

	return lam, nil
}
