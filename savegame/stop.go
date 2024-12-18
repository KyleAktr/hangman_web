package savegame

import (
	"encoding/json"
	"log"
	"io/ioutil"
)

// StopAndSaveGame : Stop the game and sore the game data in a TXT file
func StopAndSaveGame(data GameData) {
	DataGameJson, _ := json.MarshalIndent(data, "", " ")
	err := ioutil.WriteFile("save.txt", DataGameJson, 0777)
	if err != nil {
		log.Fatal(err)
	}
}