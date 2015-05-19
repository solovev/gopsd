package gopsd

import (
	"errors"
	"fmt"
)

type ImageResourcesSection struct {
	Length int32
	Blocks map[int16]*IRBlock
}

type IRBlock struct {
	Signature string
	UID       int16
	Name      string
	DataSize  int32
	Data      interface{}
}

func newImageResources() (*ImageResourcesSection, error) {
	ir := new(ImageResourcesSection)

	l, err := reader.ReadInt32()
	if err != nil {
		return nil, err
	}
	ir.Length = l
	ir.Blocks = make(map[int16]*IRBlock)

	pos := 0
	for int32(pos) < l {
		block := new(IRBlock)

		sign, err := reader.ReadString32()
		if err != nil {
			return nil, err
		}
		if sign != "8BIM" {
			return nil, errors.New(fmt.Sprintf("Wrong block [%d] signature: %s! Expected: \"8BIM\".", len(ir.Blocks), sign))
		}
		block.Signature = sign
		pos += 4

		uid, err := reader.ReadInt16()
		if err != nil {
			return nil, err
		}
		block.UID = uid
		pos += 2

		name, err := reader.ReadPascalString()
		if err != nil {
			return nil, err
		}
		block.Name = name
		pos += len(name) + 1

		size, err := reader.ReadInt32()
		if err != nil {
			return nil, err
		}
		block.DataSize = size
		pos += 4

		switch uid {
		case 1062:
			ps, err := newPrintScale()
			if err != nil {
				return nil, err
			}
			block.Data = ps
		default:
			if block.DataSize%2 != 0 {
				block.DataSize++

			}
			err = reader.Skip(block.DataSize)
			if err != nil {
				return nil, err
			}
		}
		pos += int(block.DataSize)

		ir.Blocks[uid] = block
	}

	return ir, nil
}

// UID: 1062
type PrintScale struct {
	Style int16
	X, Y  float32
	Scale float32
}

func newPrintScale() (*PrintScale, error) {
	ps := new(PrintScale)

	style, err := reader.ReadInt16()
	if err != nil {
		return nil, err
	}
	if style < 0 || style > 2 {
		return nil, errors.New(fmt.Sprintf("Wrong style of print scale resource: %d!", style))
	}
	ps.Style = style

	x, err := reader.ReadFloat32()
	if err != nil {
		return nil, err
	}
	ps.X = x

	y, err := reader.ReadFloat32()
	if err != nil {
		return nil, err
	}
	ps.Y = y

	scale, err := reader.ReadFloat32()
	if err != nil {
		return nil, err
	}
	ps.Scale = scale

	return ps, nil
}
