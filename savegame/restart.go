package savegame

import (
	"encoding/json"
	"log"
	"io/ioutil"
)

// GameData : The structure of our game
type GameData struct {
    Word             string   
    ToFind           string   
    Index            []int    
    Attempts         int      
    Tries            []string 
    LetterCheck      bool     
    HangmanPositions int 
	HangmanGraphic   string   // Nouvelle propriété pour l'étape actuelle     
}

// StartWithFlag : Restart the game with the saved data from a txt file
func StartWithFlag(Start string) GameData {
	var data GameData
	JsonData, _ := ioutil.ReadFile(Start)
	err := json.Unmarshal(JsonData, &data)
	if err != nil {
		log.Fatal(err)
	}
	return data  
}


