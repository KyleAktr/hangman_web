package game

import (
	"bufio"
	"fmt"
	"hang_web/asciifunc"
	"hang_web/position"
	"hang_web/savegame"
	"os"
	"strings"
)

// convertToSaveGameData : Explicit conversion of GameData (from game package) to GameData (from savegame package)
func convertToSaveGameData(data savegame.GameData) savegame.GameData {
	return savegame.GameData{
		Word:             data.Word,
		ToFind:           data.ToFind,
		Index:            data.Index,
		Attempts:         data.Attempts,
		Tries:            data.Tries,
		LetterCheck:      data.LetterCheck,
		HangmanPositions: data.HangmanPositions,
	}
}

// Run : Main function that handles the game loop
func Run(arg string, Data savegame.GameData) {
	Data.Attempts = 10
	reader := bufio.NewReader(os.Stdin)
	for Data.Attempts > 0 {
		fmt.Print("Choose:")
		guess, _ := reader.ReadString('\n')
		guess = strings.TrimSpace(strings.ToUpper(guess))
		for IntputTesting(guess, Data) == false {
			fmt.Print("Choose:")
			guess, _ = reader.ReadString('\n')
			guess = strings.TrimSpace(strings.ToUpper(guess))
		}
		Data.Tries = append(Data.Tries, guess)
		if len(guess) < 2 {
			Data = FindLetter(Data)
			if Data.LetterCheck == false {
				Data.Attempts--
				fmt.Print("\033[35m", "Not present in the word ", "\033[36m")
				fmt.Printf("%v attempts remaining\n", Data.Attempts)
				fmt.Print("\033[0m") // color reset
				fmt.Println(position.OpenJose()[9-Data.Attempts])
				fmt.Println(asciifunc.ToAsciiArt(Data.Word))
			} else {
				Data.Word = RevealLetters(Data)
				fmt.Println("\033[32m", "Good guess !", "\033[0m")
				fmt.Println(asciifunc.ToAsciiArt(Data.Word))
				if WordGuessed(Data) {
					fmt.Println("\033[92m", "Congrats ! You found the word.", "\033[0m")
					break
				}
			}
		} else { // if we guess a word instead of a letter or if we want to save and quit the game
			if guess == "STOP" {
				saveData := convertToSaveGameData(Data) // Convertir avant de sauvegarder
				savegame.StopAndSaveGame(saveData)
				fmt.Println("\033[32m", "Your game is saved. Enter --startWith save.txt in argument to restart your last save", "\033[0m")
				return
			}
			if GuessingWord(guess, Data) == false {
				Data.Attempts = Data.Attempts - 2
				if Data.Attempts > 0 {
					fmt.Print("\033[36m")
					fmt.Printf("You have %v attempts left\n", Data.Attempts)
					fmt.Print("\033[0m")
					fmt.Println(position.OpenJose()[9-Data.Attempts])
					fmt.Println(asciifunc.ToAsciiArt(Data.Word))
				} else {
					fmt.Println(position.OpenJose()[9])
				}
			} else {
				Data.Word = guess
				fmt.Println(asciifunc.ToAsciiArt(Data.Word))
				fmt.Println("\033[92m", "Congrats ! You found the word.", "\033[0m")
				break
			}
		}
	}
	if Data.Attempts <= 0 {
		fmt.Println("\033[31m", "Sorry, you lose! The word was:", strings.ToUpper(Data.ToFind), "\033[0m")
		if askForReplay() {
			Choosefile(arg)
		} else {
			fmt.Println("Thanks for playing! Goodbye!")
		}
	}
}

// askForReplay demande au joueur s'il veut rejouer
func askForReplay() bool {
	var response string
	fmt.Println("Do you want to play again? (yes/no):")
	fmt.Scanf("%s", &response) // On récupère la réponse de l'utilisateur

	response = strings.ToLower(response) // Convertir en minuscule pour simplifier la comparaison
	if response == "yes" || response == "y" {
		return true
	}
	return false
}

// NewGame : Start a new game
func NewGame(WordToFind string) savegame.GameData {
	return savegame.GameData{
		Word:             "",
		ToFind:           WordToFind,
		Index:            []int{},
		Attempts:         10,
		Tries:            []string{},
		LetterCheck:      false,
		HangmanPositions: 0,
	}
}

// StartGame : Start a new game
func StartGame(WordsSlice []string) savegame.GameData {
	fmt.Println("\n")
	fmt.Println(asciifunc.ToAsciiArt("HANGMAN"))
	fmt.Println()
	fmt.Println("\033[32m", "Good Luck, you have 10 attempts.\n\n", "\033[0m") // green - color reset
	WordToFind := TakeRandomWord(WordsSlice)
	Data := NewGame(WordToFind)
	Data.ToFind = WordToFind
	Data.Tries = InitialLetters(Data).Tries
	Data.Index = RevealInitialLetters(Data).Index
	Data = RevealInitialLetters(Data)
	return Data
}
