package gopsd

func readColorMode(doc *Document) {
	length := reader.ReadInt32()
	if doc.ColorMode == "Indexed" {
		// TODO
	} else if doc.ColorMode == "Duotone" {
		// TODO
	}
	reader.Skip(length)
}
