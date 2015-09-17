## gopsd
> Photoshop document parser in Golang

### Example

```go
package main

import (
	"fmt"
	"image/png"
	"os"

	"github.com/solovev/gopsd"
)

func main() {
	doc, err := gopsd.ParseFromPath("./test.psd")
	if err != nil {
		fmt.Println(err)
	}
	os.Mkdir("./images", 0777)

	for _, layer := range doc.Layers {
		fmt.Println(layer.ToString())

		saveAsPNG(layer)
	}
}

func saveAsPNG(layer *gopsd.Layer) {
	out, err := os.Create("./images/" + layer.Name + ".png")
	if err != nil {
		fmt.Println(err)
	}

	err = png.Encode(out, layer.Image)
	if err != nil {
		fmt.Println(err)
	}
}
```
### test.psd
> ![photoshop](https://raw.githubusercontent.com/solovev/gopsd/master/examples/images/readme_preview.png)

### Result

> Background: [X: 0, Y: 0, Width: 384, Height: 512] <br>
GrayRect: [X: -1, Y: 153, Width: 349, Height: 145] <br>
RoundPinkRect: [X: 66, Y: 163, Width: 273, Height: 124] <br>
Ellipse: [X: -58, Y: -89, Width: 338, Height: 212] <br>
SOME TEXT: [X: 18, Y: 18, Width: 235, Height: 34] <br>


| Layers |
| ------------- |
| Background |
| <img src="https://raw.githubusercontent.com/solovev/gopsd/master/examples/images/Background.png" alt="Background" width="192" height="256" /> |
| GrayRect |
| <img src="https://raw.githubusercontent.com/solovev/gopsd/master/examples/images/GrayRect.png" alt="GrayRect" width="175" height="72" /> |
| RoundPinkRect |
| <img src="https://raw.githubusercontent.com/solovev/gopsd/master/examples/images/RoundPinkRect.png" alt="RoundPinkRect" width="136" height="62" /> |
| Ellipse |
| <img src="https://raw.githubusercontent.com/solovev/gopsd/master/examples/images/Ellipse.png" alt="Ellipse" width="169" height="106" /> |
| SOME TEXT |
| <img src="https://raw.githubusercontent.com/solovev/gopsd/master/examples/images/SOME TEXT.png" alt="SOME TEXT" width="117" height="17" /> |
