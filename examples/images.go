package main

import (
	"fmt"
	"image/png"
	"os"

	"github.com/solovev/gopsd"
)

func main() {
	doc, err := gopsd.ParseFromPath("./test.psd")
	checkError(err)

	os.Mkdir("./images", 0777)

	for _, layer := range doc.Layers {
		fmt.Println(layer.ToString())

		saveAsPNG(layer)
	}
}

func saveAsPNG(layer *gopsd.Layer) {
	out, err := os.Create("./images/" + layer.Name + ".png")
	checkError(err)

	img, err := layer.GetImage()
	checkError(err)

	err = png.Encode(out, img)
	checkError(err)
}

func checkError(err error) {
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
