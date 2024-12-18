package asciifunc

import (
	"io/ioutil"
	"strings"
	"log"
)


// OpenAscii : This function reads the ASCII art characters from the file "standard.txt"
func OpenAscii() [8][27]string {
	var AsciiTable [8][27]string
	AsciiFile, err := ioutil.ReadFile("displaytxt/standard.txt")
	if err != nil {
		log.Fatal(err)
	}
	AsciiStr := string(AsciiFile)
	AsciiSlice := strings.Split(AsciiStr, "\n")
	x := 1
	y := 0
	for x <= 27 {
		i := 0
		for y < x*8 {
			AsciiTable[i][x-1] = AsciiSlice[y]
			y++
			i++
		}
		x++
	}
	return AsciiTable
}