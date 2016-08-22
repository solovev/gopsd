package gopsd

import (
	"image"
	"image/color"

	"github.com/solovev/gopsd/util"
)

func readImageData(doc *Document) {
	rle := reader.ReadInt16() == 1

	width := int(doc.Width)
	height := int(doc.Height)
	channels := int(doc.Channels)

	byteCounts := make([]int16, channels*height)
	if rle {
		for i := range byteCounts {
			byteCounts[i] = reader.ReadInt16()
		}
	}

	chanData := make(map[int][]int8)
	for i := 0; i < channels; i++ {
		var data []int8
		if rle {
			index := i * height
			for j := 0; j < height; j++ {
				length := byteCounts[index]
				index++

				line := util.UnpackRLEBits(reader.ReadSignedBytes(length), width)
				data = append(data, line...)
			}
		} else {
			data = reader.ReadSignedBytes(width * height)
		}
		chanData[i] = data
	}

	image := image.NewRGBA(image.Rect(0, 0, width, height))
	for x := 0; x < width; x++ {
		for y := 0; y < height; y++ {
			i := x + (y * width)
			red := byte(chanData[0][i])
			green := byte(chanData[1][i])
			blue := byte(chanData[2][i])
			image.Set(x, y, color.RGBA{red, green, blue, 255})
		}
	}
	doc.Image = image
}
