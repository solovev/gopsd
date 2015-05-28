package gopsd

import "fmt"

func readImageData(doc *Document) {
	compression := reader.ReadInt16()
	switch compression {
	case 1: // RLE
		byteCounts := make([]int16, int(doc.Channels)*int(doc.Height))
		for i := range byteCounts {
			byteCounts[i] = reader.ReadInt16()
		}

		chanData := make(map[int][]int8)
		for i := 0; i < int(doc.Channels); i++ {
			chanId := 0
			if i == 3 {
				chanId = -i
			} else {
				chanId = i
			}

			data := make([]int8, 0)
			index := i * int(doc.Height)
			for j := 0; j < int(doc.Height); j++ {
				length := byteCounts[index]
				index++
				line := readRLE(reader.ReadSignedBytes(length), int(doc.Width))
				data = append(data, line...)
			}
			fmt.Println(len(data), data[0], data[len(data)-1])
			chanData[chanId] = data
		}
	default:
		panic("Unknown compression of image data.")
	}

}

func readRLE(data []int8, length int) []int8 {
	result := make([]int8, length)
	wPos, rPos := 0, 0
	for rPos < len(data) {
		n := data[rPos]
		rPos++
		if n > 0 {
			count := int(n) + 1
			for j := 0; j < count; j++ {
				result[wPos] = data[rPos]
				wPos++
				rPos++
			}
		} else {
			b := data[rPos]
			rPos++
			count := int(-n) + 1
			for j := 0; j < count; j++ {
				result[wPos] = b
				wPos++
			}
		}
	}
	return result
}
