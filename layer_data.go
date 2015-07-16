package gopsd

import (
	"fmt"

	"github.com/solovev/gopsd/util"
)

type LayerSection struct {
	Type      int32
	SubType   int32
	BlendMode string
}

func ReadLayerSection(reader *util.Reader, length int64, layer string) *LayerSection {
	lsct := new(LayerSection)
	lsct.Type = reader.ReadInt32()
	if length >= 12 {
		sign := reader.ReadString(4)
		if sign != "8BIM" {
			panic(fmt.Sprintf("Wrong section signature of layer %s", layer))
		}
		key := reader.ReadString(4)
		if mode, ok := util.BlendModeKeys[key]; ok {
			lsct.BlendMode = mode
		}
	}
	if length >= 16 {
		lsct.SubType = reader.ReadInt32()
	}
	return lsct
}
