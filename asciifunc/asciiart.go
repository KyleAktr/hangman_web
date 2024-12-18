package asciifunc

// ToAsciiArt : Converts a word into ASCII art representation using the ASCII table from OpenAscii()
func ToAsciiArt(word string) string {
	asciiWord := ""
	AsciiTable := OpenAscii()
	var TableIndex []int
	for _, v := range word {
		if v == 95 {
			TableIndex = append(TableIndex, 0)
		} else {
			TableIndex = append(TableIndex, int(v-64))
		}
	}
	asciiWord += "\n"
	for i := 0; i < 8; i++ {
		for _, index := range TableIndex {
			asciiWord = asciiWord + "░░" + AsciiTable[i][index]
		}
		asciiWord += "░░"
		asciiWord += "\n"
	}
	return asciiWord
}