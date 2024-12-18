package position

import (
	"io/ioutil"
	"log"
	"strings"
)

// GameData : The structure of our game
type GameData struct {
	Word             string   
	ToFind           string   
	Index            []int    
	Attempts         int      
	Tries            []string 
	HangmanPositions int 
	HangmanGraphic   string   // Nouvelle propriété pour l'étape actuelle      
}

// OpenJose : Open hangman.txt file and return a string sliced fore each positions of the hangman
func OpenJose() []string {
	JoseFile, err := ioutil.ReadFile("displaytxt/hangman.txt")
	if err != nil {
		log.Fatal(err)
	}
	JoseStr := string(JoseFile)                 
	JoseSlice := strings.Split(JoseStr, "\n\n") 
	return JoseSlice
}