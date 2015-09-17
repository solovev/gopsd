## Example:
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

---
### Result
| Ellipse |
| ------------- |
| <img src="https://raw.githubusercontent.com/solovev/gopsd/master/examples/images/Ellipse.png" alt="Ellipse" width="169" height="106" /> |
